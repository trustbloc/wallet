/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package oidc

import (
	"errors"
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/crypto"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

type zcapRemoteKMS struct{}

func (r *zcapRemoteKMS) Create(kms.KeyType) (string, interface{}, error) {
	return "", nil, errors.New("zcapRemoteKMS: Create() not implemented")
}

func (r *zcapRemoteKMS) Get(keyID string) (interface{}, error) {
	return keyID, nil
}

func (r *zcapRemoteKMS) Rotate(kms.KeyType, string) (string, interface{}, error) {
	return "", nil, errors.New("zcapRemoteKMS: Rotate() not implemented")
}

func (r *zcapRemoteKMS) ExportPubKeyBytes(string) ([]byte, kms.KeyType, error) {
	return nil, "", errors.New("zcapRemoteKMS: ExportPubKeyBytes() not implemented")
}

func (r *zcapRemoteKMS) CreateAndExportPubKeyBytes(kms.KeyType) (string, []byte, error) {
	return "", nil, errors.New("zcapRemoteKMS: CreateAndExportPubKeyBytes() not implemented")
}

func (r *zcapRemoteKMS) PubKeyBytesToHandle([]byte, kms.KeyType) (interface{}, error) {
	return nil, errors.New("zcapRemoteKMS: PubKeyBytesToHandle() not implemented")
}

func (r *zcapRemoteKMS) ImportPrivateKey(interface{}, kms.KeyType, ...kms.PrivateKeyOpts) (string, interface{}, error) {
	return "", nil, errors.New("zcapRemoteKMS: ImportPrivateKey() not implemented")
}

type zcapRemoteCrypto struct {
	signer signer
}

func (r *zcapRemoteCrypto) Encrypt(_, _ []byte, _ interface{}) ([]byte, []byte, error) {
	return nil, nil, errors.New("zcapRemoteCrypto: Encrypt() not implemented")
}

func (r *zcapRemoteCrypto) Decrypt(_, _, _ []byte, _ interface{}) ([]byte, error) {
	return nil, errors.New("zcapRemoteCrypto: Decrypt() not implemented")
}

func (r *zcapRemoteCrypto) Sign(msg []byte, _ interface{}) ([]byte, error) {
	sig, err := r.signer.Sign(msg)
	if err != nil {
		return nil, fmt.Errorf("zcapRemoteCrypto: signer failed to sign: %w", err)
	}

	return sig, nil
}

func (r *zcapRemoteCrypto) Verify(_, _ []byte, _ interface{}) error {
	return errors.New("zcapRemoteCrypto: Verify() not implemented")
}

func (r *zcapRemoteCrypto) ComputeMAC([]byte, interface{}) ([]byte, error) {
	return nil, errors.New("zcapRemoteCrypto: ComputeMAC() not implemented")
}

func (r *zcapRemoteCrypto) VerifyMAC(_, _ []byte, _ interface{}) error {
	return errors.New("zcapRemoteCrypto: VerifyMAC() not implemented")
}

func (r *zcapRemoteCrypto) WrapKey(_, _, _ []byte,
	_ *crypto.PublicKey, _ ...crypto.WrapKeyOpts) (*crypto.RecipientWrappedKey, error) {
	return nil, errors.New("zcapRemoteCrypto: WrapKey() not implemented")
}

func (r *zcapRemoteCrypto) UnwrapKey(*crypto.RecipientWrappedKey, interface{}, ...crypto.WrapKeyOpts) ([]byte, error) {
	return nil, errors.New("zcapRemoteCrypto: UnwrapKey() not implemented")
}

func (r *zcapRemoteCrypto) SignMulti(messages [][]byte, kh interface{}) ([]byte, error) {
	return nil, errors.New("zcapRemoteCrypto: SignMulti() not implemented")
}

func (r *zcapRemoteCrypto) VerifyMulti(messages [][]byte, signature []byte, kh interface{}) error {
	return errors.New("zcapRemoteCrypto: VerifyMulti() not implemented")
}

func (r *zcapRemoteCrypto) VerifyProof(revealedMessages [][]byte, proof, nonce []byte, kh interface{}) error {
	return errors.New("zcapRemoteCrypto: VerifyProof() not implemented")
}

func (r *zcapRemoteCrypto) DeriveProof(messages [][]byte, bbsSignature, nonce []byte, revealedIndexes []int,
	kh interface{}) ([]byte, error) {
	return nil, errors.New("zcapRemoteCrypto: DeriveProof() not implemented")
}
