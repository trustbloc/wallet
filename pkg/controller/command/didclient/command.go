/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package didclient

import (
	"fmt"
	"io"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	ariesapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/hyperledger/aries-framework-go/pkg/kms/legacykms"
	ariesstorage "github.com/hyperledger/aries-framework-go/pkg/storage"
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
	// CreateDIDErrorCode is typically a code for validation errors
	// for create did error
	CreateDIDErrorCode = command.Code(iota + command.DIDClient)
)

type didBlocClient interface {
	CreateDID(domain string) (*did.Doc, error)
}

// New returns new DID Exchange controller command instance
func New(storeProvider ariesstorage.Provider, domain string) (*Command, error) {
	// Create KMS
	kms, err := createKMS(storeProvider)
	if err != nil {
		return nil, err
	}

	client := didclient.New(kms)

	cmd := &Command{
		client: client,
		domain: domain,
	}

	return cmd, nil
}

func createKMS(s ariesstorage.Provider) (ariesapi.CloseableKMS, error) {
	kmsProvider, err := context.New(context.WithStorageProvider(s))
	if err != nil {
		return nil, fmt.Errorf("failed to create new kms provider: %w", err)
	}

	kms, err := legacykms.New(kmsProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new kms: %w", err)
	}

	return kms, nil
}

// Command is controller command for DID Exchange
type Command struct {
	client didBlocClient
	domain string
}

// GetHandlers returns list of all commands supported by this controller command
func (c *Command) GetHandlers() []command.Handler {
	return []command.Handler{
		cmdutil.NewCommandHandler(commandName, createDIDCommandMethod, c.CreateDID),
	}
}

// CreateInvitation Creates a new connection invitation.
func (c *Command) CreateDID(rw io.Writer, req io.Reader) command.Error {
	didDoc, err := c.client.CreateDID(c.domain)
	if err != nil {
		logutil.LogError(logger, commandName, createDIDCommandMethod, err.Error())
		return command.NewExecuteError(CreateDIDErrorCode, err)
	}

	command.WriteNillableResponse(rw, didDoc, logger)

	logutil.LogDebug(logger, commandName, createDIDCommandMethod, successString)

	return nil
}
