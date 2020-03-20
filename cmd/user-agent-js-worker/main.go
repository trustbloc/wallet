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
	"syscall/js"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/storage/jsindexeddb"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/trustbloc/edge-agent/pkg/controller"
	cmdctrl "github.com/trustbloc/edge-agent/pkg/controller/command"
)

var logger = logrus.New()

const (
	wasmStartupTopic    = "asset-ready"
	handleResultFn      = "handleResult"
	userAgentCommandPkg = "useragent"
	userAgentStartFn    = "Start"
	userAgentStopFn     = "Stop"
)

// TODO Signal JS when WASM is loaded and ready.
//      This is being used in tests for now.
var ready = make(chan struct{}) //nolint:gochecknoglobals
var isTest = false              //nolint:gochecknoglobals

// command is received from JS
type command struct {
	ID      string                 `json:"id"`
	Pkg     string                 `json:"pkg"`
	Fn      string                 `json:"fn"`
	Payload map[string]interface{} `json:"payload"`
}

// result is sent back to JS
type result struct {
	ID      string                 `json:"id"`
	IsErr   bool                   `json:"isErr"`
	ErrMsg  string                 `json:"errMsg"`
	Payload map[string]interface{} `json:"payload"`
	Topic   string                 `json:"topic"`
}

// userAgentStartOpts contains opts for starting user agent
type userAgentStartOpts struct {
	BlocDomain string `json:"blocDomain"`
}

// main registers the 'handleMsg' function in the JS context's global scope to receive commands.
// results are posted back to the 'handleResult' JS function.
func main() {
	input := make(chan *command)
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
			logger.Warnf("useragent wasm: unable to unmarshal input=%s. err=%s", args[0].String(), err)

			return nil
		}
		in <- cmd

		return nil
	}
}

func pipe(input chan *command, output chan *result) {
	handlers := testHandlers()

	addUserAgentHandlers(handlers)

	for c := range input {
		if c.ID == "" {
			logger.Warnf("useragent wasm: missing ID for input: %v", c)
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
			logger.Errorf("useragent wasm: failed to marshal response for id=%s err=%s ", r.ID, err)
		}

		js.Global().Call(handleResultFn, string(out))
	}
}

func cmdExecToFn(exec cmdctrl.Exec) func(*command) *result {
	return func(c *command) *result {
		b, er := json.Marshal(c.Payload)
		if er != nil {
			return &result{
				ID:     c.ID,
				IsErr:  true,
				ErrMsg: fmt.Sprintf("user wasm: failed to unmarshal payload. err=%s", er),
			}
		}

		req := bytes.NewBuffer(b)

		var buf bytes.Buffer

		err := exec(&buf, req)
		if err != nil {
			return newErrResult(c.ID, fmt.Sprintf("code: %+v, message: %s", err.Code(), err.Error()))
		}

		payload := make(map[string]interface{})

		if len(buf.Bytes()) > 0 {
			if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
				return &result{
					ID:    c.ID,
					IsErr: true,
					ErrMsg: fmt.Sprintf(
						"user wasm: failed to unmarshal command result=%+v err=%s",
						buf.String(), err),
				}
			}
		}

		return &result{
			ID:      c.ID,
			Payload: payload,
		}
	}
}

func addUserAgentHandlers(pkgMap map[string]map[string]func(*command) *result) {
	fnMap := make(map[string]func(*command) *result)
	fnMap[userAgentStartFn] = func(c *command) *result {
		cOpts, err := startOpts(c.Payload)
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		s, _ := jsindexeddb.NewProvider("")
		commands, err := controller.GetCommandHandlers(s, controller.WithBlocDomain(cOpts.BlocDomain))
		if err != nil {
			return newErrResult(c.ID, err.Error())
		}

		// add command handlers
		addCommandHandlers(commands, pkgMap)

		// add stop user agent handler
		addStopUserAgentHandler(pkgMap)

		return &result{
			ID:      c.ID,
			Payload: map[string]interface{}{"message": "user agent started successfully"},
		}
	}

	pkgMap[userAgentCommandPkg] = fnMap
}

func addCommandHandlers(commands []cmdctrl.Handler, pkgMap map[string]map[string]func(*command) *result) {
	for _, cmd := range commands {
		fnMap, ok := pkgMap[cmd.Name()]
		if !ok {
			fnMap = make(map[string]func(*command) *result)
		}

		fnMap[cmd.Method()] = cmdExecToFn(cmd.Handle())
		pkgMap[cmd.Name()] = fnMap
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

func addStopUserAgentHandler(pkgMap map[string]map[string]func(*command) *result) {
	fnMap := make(map[string]func(*command) *result)
	fnMap[userAgentStopFn] = func(c *command) *result {
		// reset handlers when stopped
		for k := range pkgMap {
			delete(pkgMap, k)
		}

		// put back start command once stopped
		addUserAgentHandlers(pkgMap)

		return &result{
			ID:      c.ID,
			Payload: map[string]interface{}{"message": "user agent stopped"},
		}
	}
	pkgMap[userAgentCommandPkg] = fnMap
}

func isStartCommand(c *command) bool {
	return c.Pkg == userAgentCommandPkg && c.Fn == userAgentStartFn
}

func isStopCommand(c *command) bool {
	return c.Pkg == userAgentCommandPkg && c.Fn == userAgentStopFn
}

func handlerNotFoundErr(c *command) *result {
	if isStartCommand(c) {
		return newErrResult(c.ID, "user agent already started")
	} else if isStopCommand(c) {
		return newErrResult(c.ID, "user agent not running")
	}

	return newErrResult(c.ID, fmt.Sprintf("invalid pkg/fn: %s/%s, make sure user agent is started", c.Pkg, c.Fn))
}

func newErrResult(id, msg string) *result {
	return &result{
		ID:     id,
		IsErr:  true,
		ErrMsg: "useragent wasm: " + msg,
	}
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

func startOpts(payload map[string]interface{}) (*userAgentStartOpts, error) {
	opts := &userAgentStartOpts{}

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
