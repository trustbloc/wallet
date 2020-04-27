/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package didclient

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"io"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/sirupsen/logrus"
	didclient "github.com/trustbloc/trustbloc-did-method/pkg/did"

	"github.com/trustbloc/edge-agent/pkg/controller/command"
	"github.com/trustbloc/edge-agent/pkg/controller/command/internal/cmdutil"
	"github.com/trustbloc/edge-agent/pkg/controller/command/internal/logutil"
)

var logger = logrus.New()

const (
	// command name
	commandName = "didclient"
	// command methods
	createDIDCommandMethod = "CreateDID"
	// log constants
	successString = "success"
)

const (
	// InvalidRequestErrorCode is typically a code for validation errors
	InvalidRequestErrorCode = command.Code(iota + command.DIDClient)

	// CreateDIDErrorCode is typically a code for create did errors
	CreateDIDErrorCode

	// GenerateKeyPairErrorCode is an error code for indicating key pair failure
	GenerateKeyPairErrorCode
)

type didBlocClient interface {
	CreateDID(domain string, opts ...didclient.CreateDIDOption) (*did.Doc, error)
}

// New returns new DID Exchange controller command instance
func New(domain string) *Command {
	client := didclient.New()

	cmd := &Command{
		client: client,
		domain: domain,
	}

	cmd.generateECKeyPair = func() ([]byte, []byte, error) {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		publicKeyBytes := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
		// TODO do we need to use MarshalECPrivateKey
		encodedPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return nil, nil, err
		}
		return publicKeyBytes, encodedPrivateKey, nil
	}
	return cmd
}

// Command is controller command for DID Exchange
type Command struct {
	client            didBlocClient
	domain            string
	generateECKeyPair func() ([]byte, []byte, error) // needed for unit test
}

// GetHandlers returns list of all commands supported by this controller command
func (c *Command) GetHandlers() []command.Handler {
	return []command.Handler{
		cmdutil.NewCommandHandler(commandName, createDIDCommandMethod, c.CreateDID),
	}
}

// CreateInvitation Creates a new connection invitation.
func (c *Command) CreateDID(rw io.Writer, req io.Reader) command.Error {
	var request CreateDIDRequest

	err := json.NewDecoder(req).Decode(&request)
	if err != nil {
		logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
		return command.NewValidationError(InvalidRequestErrorCode, err)
	}

	var opts []didclient.CreateDIDOption
	var didPrivateKey string

	for _, v := range request.PublicKeys {
		switch v.KeyType {
		case didclient.Ed25519KeyType:
			opts = append(opts, didclient.WithPublicKey(&didclient.PublicKey{ID: v.ID, Type: v.Type, Encoding: v.Encoding,
				KeyType: v.KeyType, Usage: v.Usage, Recovery: v.Recovery, Value: base58.Decode(v.Value)}))
		case didclient.P256KeyType:
			encodedPublicKey, encodedPrivateKey, err := c.generateECKeyPair()
			if err != nil {
				logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
				return command.NewExecuteError(GenerateKeyPairErrorCode, err)
			}
			didPrivateKey = base58.Encode(encodedPrivateKey)
			opts = append(opts, didclient.WithPublicKey(&didclient.PublicKey{ID: v.ID, Type: v.Type, Encoding: v.Encoding,
				KeyType: v.KeyType, Usage: v.Usage, Recovery: v.Recovery, Value: encodedPublicKey}))
		}
	}

	didDoc, err := c.client.CreateDID(c.domain, opts...)
	if err != nil {
		logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
		return command.NewExecuteError(CreateDIDErrorCode, err)
	}

	bytes, err := didDoc.JSONBytes()
	if err != nil {
		logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
		return command.NewExecuteError(CreateDIDErrorCode, err)
	}

	m := make(map[string]interface{})

	if err := json.Unmarshal(bytes, &m); err != nil {
		logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
		return command.NewExecuteError(CreateDIDErrorCode, err)
	}

	command.WriteNillableResponse(rw, &CreateDIDResponse{
		DID:        m,
		PrivateKey: didPrivateKey,
	}, logger)

	logutil.LogDebug(logger, commandName, createDIDCommandMethod, successString)

	return nil
}
