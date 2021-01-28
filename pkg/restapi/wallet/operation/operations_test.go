/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package operation // nolint:testpackage // changing to different package requires exposing internal features.

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("create new instance - success", func(t *testing.T) {
		op, err := New(nil, nil, nil)

		require.NoError(t, err)
		require.Len(t, op.GetRESTHandlers(), 0)
	})
}
