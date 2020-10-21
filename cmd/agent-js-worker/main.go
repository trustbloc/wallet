// +build js,wasm

/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall/js"
	"time"

	"github.com/google/uuid"

	"github.com/hyperledger/aries-framework-go/component/storage/jsindexeddb"
	ariesctrl "github.com/hyperledger/aries-framework-go/pkg/controller"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/messaging/msghandler"
	arieshttp "github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/http"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/ws"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/httpbinding"

	"github.com/mitchellh/mapstructure"
	"github.com/trustbloc/edge-core/pkg/log"

	agentctrl "github.com/trustbloc/edge-agent/pkg/controller"
)

var logger = log.New("agent-js-worker")

const (
	wasmStartupTopic = "asset-ready"
	handleResultFn   = "handleResult"
	commandPkg       = "agent"
	startFn          = "Start"
	stopFn           = "Stop"
	workers          = 2
)

// TODO Signal JS when WASM is loaded and ready.
//      This is being used in tests for now.
var (
	ready  = make(chan struct{}) //nolint:gochecknoglobals
	isTest = false               //nolint:gochecknoglobals
)

// command is received from JS.
type command struct {
	ID      string                 `json:"id"`
	Pkg     string                 `json:"pkg"`
	Fn      string                 `json:"fn"`
	Payload map[string]interface{} `json:"payload"`
}

// result is sent back to JS.
type result struct {
	ID      string                 `json:"id"`
	IsErr   bool                   `json:"isErr"`
	ErrMsg  string                 `json:"errMsg"`
	Payload map[string]interface{} `json:"payload"`
	Topic   string                 `json:"topic"`
}

// agentStartOpts contains opts for starting agent.
type agentStartOpts struct {
	Label                string   `json:"agent-default-label"`
	HTTPResolvers        []string `json:"http-resolver-url"`
	AutoAccept           bool     `json:"auto-accept"`
	OutboundTransport    []string `json:"outbound-transport"`
	TransportReturnRoute string   `json:"transport-return-route"`
	LogLevel             string   `json:"log-level"`
	DBNamespace          string   `json:"db-namespace"`
	BlocDomain           string   `json:"blocDomain"`
	SDSServerURL         string   `json:"sdsServerURL"`
}

// main registers the 'handleMsg' function in the JS context's global scope to receive commands.
// Results are posted back to the 'handleResult' JS function.
func main() {
	// TODO: capacity was added due to deadlock. Looks like js worker are not able to pick up 'output chan *result'.
	//  Another fix for that is to wrap 'in <- cmd' in a goroutine. e.g go func() { in <- cmd }()
	//  We need to figure out what is the root cause of deadlock and fix it properly.
	input := make(chan *command, 10)
	output := make(chan *result)

	go pipe(input, output)

	go sendTo(output)

	js.Global().Set("handleMsg", js.FuncOf(takeFrom(input)))

	postInitMsg()

	if isTest {
		ready <- struct{}{}
	}

	select {}
}

func takeFrom(in chan *command) func(js.Value, []js.Value) interface{} {
	return func(_ js.Value, args []js.Value) interface{} {
		cmd := &command{}
		if err := json.Unmarshal([]byte(args[0].String()), cmd); err != nil {
			logger.Errorf("agent wasm: unable to unmarshal input=%s. err=%s", args[0].String(), err)

			return nil
		}

		in <- cmd

		return nil
	}
}

func pipe(input chan *command, output chan *result) {
	handlers := testHandlers()

	addAgentHandlers(handlers)

	for w := 0; w < workers; w++ {
		go worker(input, output, handlers)
	}
}

func worker(input chan *command, output chan *result, handlers map[string]map[string]func(*command) *result) {
	for c := range input {
		if c.ID == "" {
			logger.Warnf("agent wasm: missing ID for input: %v", c)
		}

		if pkg, found := handlers[c.Pkg]; found {
			if fn, found := pkg[c.Fn]; found {
				output <- fn(c)
				continue
			}
		}

		output <- handlerNotFoundErr(c)
	}
}

func sendTo(out chan *result) {
	for r := range out {
		out, err := json.Marshal(r)
		if err != nil {
			logger.Errorf("agent wasm: failed to marshal response for id=%s err=%s ", r.ID, err)
		}

		js.Global().Call(handleResultFn, string(out))
	}
}

func testHandlers() map[string]map[string]func(*command) *result {
	return map[string]map[string]func(*command) *result{
		"test": {
			"echo": func(c *command) *result {
				return &result{
					ID:      c.ID,
					Payload: map[string]interface{}{"echo": c.Payload},
				}
			},
			"throwError": func(c *command) *result {
				return newErrResult(c.ID, "an error !!")
			},
			"timeout": func(c *command) *result {
				const echoTimeout = 10 * time.Second

				time.Sleep(echoTimeout)

				return &result{
					ID:      c.ID,
					Payload: map[string]interface{}{"echo": c.Payload},
				}
			},
		},
	}
}

func isStartCommand(c *command) bool {
	return c.Pkg == commandPkg && c.Fn == startFn
}

func isStopCommand(c *command) bool {
	return c.Pkg == commandPkg && c.Fn == stopFn
}

func handlerNotFoundErr(c *command) *result {
	if isStartCommand(c) {
		return newErrResult(c.ID, "Agent already started")
	} else if isStopCommand(c) {
		return newErrResult(c.ID, "Agent not running")
	}

	return newErrResult(c.ID, fmt.Sprintf("invalid pkg/fn: %s/%s, make sure agent is started", c.Pkg, c.Fn))
}

func addAgentHandlers(pkgMap map[string]map[string]func(*command) *result) {
	fnMap := make(map[string]func(*command) *result)
	fnMap[startFn] = func(c *command) *result {
		cOpts, err := startOpts(c.Payload)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		err = setLogLevel(cOpts.LogLevel)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		options, err := agentOpts(cOpts)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		msgHandler := msghandler.NewRegistrar()
		options = append(options, aries.WithMessageServiceProvider(msgHandler))

		a, err := aries.New(options...)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		handlers, err := getAriesHandlers(a, msgHandler, cOpts)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		agentHandlers, err := getAgentHandlers(cOpts)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		handlers = append(handlers, agentHandlers...)

		// add command handlers
		addCommandHandlers(handlers, pkgMap)

		// add stop agent handler
		addStopAgentHandler(a, pkgMap)

		return &result{
			ID:      c.ID,
			Payload: map[string]interface{}{"message": "agent started successfully"},
		}
	}

	pkgMap[commandPkg] = fnMap
}

type execFn func(rw io.Writer, req io.Reader) error

type commandHandler struct {
	name   string
	method string
	exec   execFn
}

func getAriesHandlers(a *aries.Aries, r *msghandler.Registrar, opts *agentStartOpts) ([]commandHandler, error) {
	ctx, err := a.Context()
	if err != nil {
		return nil, err
	}

	handlers, err := ariesctrl.GetCommandHandlers(ctx, ariesctrl.WithMessageHandler(r),
		ariesctrl.WithDefaultLabel(opts.Label), ariesctrl.WithNotifier(&jsNotifier{}))
	if err != nil {
		return nil, err
	}

	var hh []commandHandler

	for _, h := range handlers {
		handle := h.Handle()
		hh = append(hh, commandHandler{
			name:   h.Name(),
			method: h.Method(),
			exec: func(rw io.Writer, req io.Reader) error {
				e := handle(rw, req)
				if e != nil {
					return fmt.Errorf("code: %+v, message: %s", e.Code(), e.Error())
				}
				return nil
			},
		})
	}

	return hh, nil
}

func getAgentHandlers(opts *agentStartOpts) ([]commandHandler, error) {
	handlers, err := agentctrl.GetCommandHandlers(agentctrl.WithBlocDomain(opts.BlocDomain),
		agentctrl.WithSDSServerURL(opts.SDSServerURL))
	if err != nil {
		return nil, err
	}

	var hh []commandHandler

	for _, h := range handlers {
		handle := h.Handle()
		hh = append(hh, commandHandler{
			name:   h.Name(),
			method: h.Method(),
			exec: func(rw io.Writer, req io.Reader) error {
				e := handle(rw, req)
				if e != nil {
					return fmt.Errorf("code: %+v, message: %s", e.Code(), e.Error())
				}
				return nil
			},
		})
	}

	return hh, nil
}

func addCommandHandlers(handlers []commandHandler, pkgMap map[string]map[string]func(*command) *result) {
	for _, h := range handlers {
		fnMap, ok := pkgMap[h.name]
		if !ok {
			fnMap = make(map[string]func(*command) *result)
		}

		fnMap[h.method] = cmdExecToFn(h.exec)
		pkgMap[h.name] = fnMap
	}
}

func cmdExecToFn(exec execFn) func(*command) *result {
	return func(c *command) *result {
		b, er := json.Marshal(c.Payload)
		if er != nil {
			return &result{
				ID:     c.ID,
				IsErr:  true,
				ErrMsg: fmt.Sprintf("agent wasm: failed to unmarshal payload. err=%s", er),
			}
		}

		req := bytes.NewBuffer(b)

		var buf bytes.Buffer

		err := exec(&buf, req)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		payload := make(map[string]interface{})

		if len(buf.Bytes()) > 0 {
			if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
				return &result{
					ID:     c.ID,
					IsErr:  true,
					ErrMsg: fmt.Sprintf("agent wasm: failed to unmarshal command result=%+v err=%s", buf.String(), err),
				}
			}
		}

		return &result{
			ID:      c.ID,
			Payload: payload,
		}
	}
}

func addStopAgentHandler(a *aries.Aries, pkgMap map[string]map[string]func(*command) *result) {
	fnMap := make(map[string]func(*command) *result)
	fnMap[stopFn] = func(c *command) *result {
		err := a.Close()
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		// reset handlers when stopped
		for k := range pkgMap {
			delete(pkgMap, k)
		}

		// put back start command once stopped
		addAgentHandlers(pkgMap)

		return &result{
			ID:      c.ID,
			Payload: map[string]interface{}{"message": "agent stopped"},
		}
	}
	pkgMap[commandPkg] = fnMap
}

func newErrResult(id, msg string) *result {
	return &result{
		ID:     id,
		IsErr:  true,
		ErrMsg: "agent wasm: " + msg,
	}
}

func startOpts(payload map[string]interface{}) (*agentStartOpts, error) {
	opts := &agentStartOpts{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  opts,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(payload)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func agentOpts(opts *agentStartOpts) ([]aries.Option, error) {
	msgHandler := msghandler.NewRegistrar()

	var options []aries.Option
	options = append(options, aries.WithMessageServiceProvider(msgHandler))

	if opts.TransportReturnRoute != "" {
		options = append(options, aries.WithTransportReturnRoute(opts.TransportReturnRoute))
	}

	store, err := jsindexeddb.NewProvider(opts.DBNamespace)
	if err != nil {
		return nil, err
	}

	options = append(options, aries.WithStoreProvider(store))

	for _, transport := range opts.OutboundTransport {
		switch transport {
		case "http":
			outbound, err := arieshttp.NewOutbound(arieshttp.WithOutboundHTTPClient(&http.Client{}))
			if err != nil {
				return nil, err
			}

			options = append(options, aries.WithOutboundTransports(outbound))
		case "ws":
			options = append(options, aries.WithOutboundTransports(ws.NewOutbound()))
		default:
			return nil, fmt.Errorf("unsupported transport : %s", transport)
		}
	}

	if len(opts.HTTPResolvers) > 0 {
		rsopts, err := getResolverOpts(opts.HTTPResolvers)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare http resolver opts : %w", err)
		}

		options = append(options, rsopts...)
	}

	return options, nil
}

func getResolverOpts(httpResolvers []string) ([]aries.Option, error) {
	var opts []aries.Option

	const numPartsResolverOption = 2

	if len(httpResolvers) > 0 {
		for _, httpResolver := range httpResolvers {
			r := strings.Split(httpResolver, "@")
			if len(r) != numPartsResolverOption {
				return nil, fmt.Errorf("invalid http resolver options found")
			}

			httpVDRI, err := httpbinding.New(r[1],
				httpbinding.WithAccept(func(method string) bool { return method == r[0] }))
			if err != nil {
				return nil, fmt.Errorf("failed to setup http resolver :  %w", err)
			}

			opts = append(opts, aries.WithVDR(httpVDRI))
		}
	}

	return opts, nil
}

func setLogLevel(logLevel string) error {
	if logLevel != "" {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			return err
		}

		log.SetLevel("", level)
		logger.Infof("log level set to `%s`", logLevel)
	}

	return nil
}

// jsNotifier notifies about all incoming events.
type jsNotifier struct {
}

// Notify is mock implementation of webhook notifier Notify().
func (n *jsNotifier) Notify(topic string, message []byte) error {
	payload := make(map[string]interface{})
	if err := json.Unmarshal(message, &payload); err != nil {
		return err
	}

	out, err := json.Marshal(&result{
		ID:      uuid.New().String(),
		Topic:   topic,
		Payload: payload,
	})
	if err != nil {
		return err
	}

	js.Global().Call(handleResultFn, string(out))

	return nil
}

func postInitMsg() {
	if isTest {
		return
	}

	out, err := json.Marshal(&result{
		ID:    uuid.New().String(),
		Topic: wasmStartupTopic,
	})
	if err != nil {
		panic(err)
	}

	js.Global().Call(handleResultFn, string(out))
}
