/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package command

import (
	"io"
)

// Exec is controller command execution function type.
type Exec func(rw io.Writer, req io.Reader) Error

// Handler for each controller command.
type Handler interface {
	// name of the command
	Name() string
	// method name of the command
	Method() string
	// execute function of the command
	Handle() Exec
}
