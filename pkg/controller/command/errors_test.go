/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewExecuteError(t *testing.T) {
	e := NewExecuteError(1, fmt.Errorf("error"))
	require.Equal(t, "error", e.Error())
	require.Equal(t, Code(1), e.Code())
	require.Equal(t, Type(0), e.Type())

}
