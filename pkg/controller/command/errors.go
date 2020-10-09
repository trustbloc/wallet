/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package command

// Type is command error type.
type Type int32

const (
	// ValidationError is error type for command validation errors.
	ValidationError Type = iota
	// ExecuteError is error type for command execution failure.
	ExecuteError Type = iota
)

// Code is the error code of command errors.
type Code int32

// Group is the error groups.
// Note: recommended to use [0-9]*000 pattern for any new entries.
// Example: 2000, 3000, 4000 ...... 25000.
type Group int32

const (
	// DIDClient error group for DID client command errors.
	DIDClient Group = 1000
)

// Error is the  interface for representing an command error condition, with the nil value representing no error.
type Error interface {
	error
	// Code returns error code for this command error.
	Code() Code
	// Code returns error type for this command error.
	Type() Type
}

// NewExecuteError returns new command execute error.
func NewExecuteError(code Code, err error) Error {
	return &commandError{err, code, ExecuteError}
}

// NewValidationError returns new command validation error.
func NewValidationError(code Code, err error) Error {
	return &commandError{err, code, ValidationError}
}

// commandError implements basic command Error.
type commandError struct {
	error
	code    Code
	errType Type
}

func (c *commandError) Code() Code {
	return c.code
}

func (c *commandError) Type() Type {
	return c.errType
}
