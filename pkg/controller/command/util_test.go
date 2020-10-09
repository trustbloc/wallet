/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package command_test

import (
	"testing"

	"github.com/trustbloc/edge-agent/pkg/controller/command"
	"github.com/trustbloc/edge-core/pkg/log"
)

func Test_WriteNillableResponse(t *testing.T) {
	command.WriteNillableResponse(&mockWriter{}, nil, log.New("util-test"))
}

type mockWriter struct {
}

func (m *mockWriter) Write(_ []byte) (n int, err error) {
	return 0, nil
}
