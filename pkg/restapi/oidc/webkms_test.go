/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc // nolint:testpackage // testing package-private types

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestZCapRemoteKMS(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, _, err := r.Create("123")
		require.EqualError(t, err, "zcapRemoteKMS: Create() not implemented")
	})

	t.Run("Get", func(t *testing.T) {
		keyID := uuid.New().String()
		r := &zcapRemoteKMS{}
		result, err := r.Get(keyID)
		require.NoError(t, err)
		require.Equal(t, keyID, result)
	})

	t.Run("Rotate", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, _, err := r.Rotate("", "")
		require.EqualError(t, err, "zcapRemoteKMS: Rotate() not implemented")
	})

	t.Run("ExportPubKeyBytes", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, _, err := r.ExportPubKeyBytes("")
		require.EqualError(t, err, "zcapRemoteKMS: ExportPubKeyBytes() not implemented")
	})

	t.Run("CreateAndExportPubKeyBytes", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, _, err := r.CreateAndExportPubKeyBytes("")
		require.EqualError(t, err, "zcapRemoteKMS: CreateAndExportPubKeyBytes() not implemented")
	})

	t.Run("PubKeyBytesToHandle", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, err := r.PubKeyBytesToHandle(nil, "")
		require.EqualError(t, err, "zcapRemoteKMS: PubKeyBytesToHandle() not implemented")
	})

	t.Run("ImportPrivateKey", func(t *testing.T) {
		r := &zcapRemoteKMS{}
		_, _, err := r.ImportPrivateKey(nil, "")
		require.EqualError(t, err, "zcapRemoteKMS: ImportPrivateKey() not implemented")
	})
}

func TestZcapRemoteCrypto(t *testing.T) {
	t.Run("Encrypt", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, _, err := r.Encrypt(nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: Encrypt() not implemented")
	})

	t.Run("Decrypt", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.Decrypt(nil, nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: Decrypt() not implemented")
	})

	t.Run("Sign", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			expected := uuid.New().String()
			r := &zcapRemoteCrypto{signer: &mockSigner{signVal: []byte(expected)}}
			result, err := r.Sign([]byte("msg"), nil)
			require.NoError(t, err)
			require.Equal(t, expected, string(result))
		})

		t.Run("error", func(t *testing.T) {
			expected := errors.New("test")
			r := &zcapRemoteCrypto{signer: &mockSigner{signErr: expected}}
			_, err := r.Sign([]byte("msg"), nil)
			require.Error(t, err)
			require.True(t, errors.Is(err, expected))
		})
	})

	t.Run("Verify", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		err := r.Verify(nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: Verify() not implemented")
	})

	t.Run("ComputeMAC", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.ComputeMAC(nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: ComputeMAC() not implemented")
	})

	t.Run("VerifyMAC", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		err := r.VerifyMAC(nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: VerifyMAC() not implemented")
	})

	t.Run("WrapKey", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.WrapKey(nil, nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: WrapKey() not implemented")
	})

	t.Run("UnwrapKey", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.UnwrapKey(nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: UnwrapKey() not implemented")
	})

	t.Run("SignMulti", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.SignMulti(nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: SignMulti() not implemented")
	})

	t.Run("VerifyMulti", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		err := r.VerifyMulti(nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: VerifyMulti() not implemented")
	})

	t.Run("VerifyProof", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		err := r.VerifyProof(nil, nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: VerifyProof() not implemented")
	})

	t.Run("DeriveProof", func(t *testing.T) {
		r := &zcapRemoteCrypto{}
		_, err := r.DeriveProof(nil, nil, nil, nil, nil)
		require.EqualError(t, err, "zcapRemoteCrypto: DeriveProof() not implemented")
	})
}
