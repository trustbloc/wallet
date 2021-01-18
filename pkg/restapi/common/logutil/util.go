/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package logutil provides helpful functions for logging.
package logutil

import (
	"github.com/trustbloc/edge-core/pkg/log"
)

// LogError is a utility function to log error messages.
func LogError(logger log.Logger, command, action, errMsg string, data ...string) {
	logger.Errorf("command=[%s] action=[%s] %s errMsg=[%s]", command, action, data, errMsg)
}

// LogDebug is a utility function to log debug messages.
func LogDebug(logger log.Logger, command, action, msg string, data ...string) {
	logger.Debugf("command=[%s] action=[%s] %s msg=[%s]", command, action, data, msg)
}

// LogInfo is a utility function to log info messages.
func LogInfo(logger log.Logger, command, action, msg string, data ...string) {
	logger.Infof("command=[%s] action=[%s] %s msg=[%s]", command, action, data, msg)
}
