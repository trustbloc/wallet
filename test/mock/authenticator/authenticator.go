/*
Copyright SecureKey Technologies Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package authenticator

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"path"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/protocol/googletpm"
	"github.com/duo-labs/webauthn/protocol/webauthncose"
	"github.com/fxamacker/cbor/v2"
)

type responseData struct {
	AttestationObject string `json:"attestationObject"`
	ClientDataJSON    string `json:"clientDataJSON"`
}

type credentialResponse struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Response responseData `json:"response"`
}

const (
	userPresentFlag  = 1
	userVerifiedFlag = 4
	hasAttdata       = 64
)

type pubKeyCredSource struct {
	Type   string
	ID     string
	Priv   *ecdsa.PrivateKey
	RPID   string
	userID []byte
}

// MockAuthenticator is a functional mock of a webauthn/FIDO authenticator device.
//  It supports fido-u2f attestation and assertion, using the P256 curve, with direct attestation,
//  using a provided CA certificate and private key to serve as an Attestation CA.
type MockAuthenticator struct {
	CredSources map[string]*pubKeyCredSource
}

// New creates a MockAuthenticator.
func New() *MockAuthenticator {
	return &MockAuthenticator{
		CredSources: map[string]*pubKeyCredSource{},
	}
}

// Authenticate generates a webauthn credential based on the provided parameters, signed by the given parent CA.
// Returns a JSON object containing the authentication response.
func (ma *MockAuthenticator) Authenticate(parentCertPath, parentKeyPath, requesterOrigin string, params *protocol.CredentialCreation) ([]byte, error) { // nolint:lll,funlen // it's just a long procedure
	clientDataJSON, clientDataHash, err := collectClientData(requesterOrigin, params.Response.Challenge)
	if err != nil {
		return nil, err
	}

	parentCert, err := loadCert(parentCertPath)
	if err != nil {
		return nil, err
	}

	parentKey, err := loadKey(parentKeyPath)
	if err != nil {
		return nil, err
	}

	myPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	myPubCose := webauthncose.EC2PublicKeyData{
		PublicKeyData: webauthncose.PublicKeyData{
			KeyType:   int64(webauthncose.EllipticKey),
			Algorithm: int64(webauthncose.AlgES256), // mock authenticator is hardcoded to use ES256
		},
		Curve:  int64(googletpm.CurveNISTP256),
		XCoord: myPriv.X.Bytes(),
		YCoord: myPriv.Y.Bytes(),
	}

	myPubCBOR, err := cbor.Marshal(myPubCose)
	if err != nil {
		return nil, err
	}

	randCredID := make([]byte, 16)

	if _, e := rand.Reader.Read(randCredID); e != nil {
		return nil, e
	}

	rpIDHash := sha256.Sum256([]byte(params.Response.RelyingParty.ID))

	authData := protocol.AuthenticatorData{
		RPIDHash: rpIDHash[:],
		Flags:    userPresentFlag | userVerifiedFlag | hasAttdata,
		AttData: protocol.AttestedCredentialData{
			AAGUID:              bytes.Repeat([]byte{0}, 16), // fido-u2f AAGUID is 16 zeros
			CredentialID:        randCredID,
			CredentialPublicKey: myPubCBOR, // pub key marshaled in webauthn COSE CBOR format
		},
	}

	authDataBytes := marshalAuthData(&authData, true)

	// generate contents of AttStatement
	sig, err := makeSig(rpIDHash[:], clientDataHash, randCredID, myPriv)
	if err != nil {
		return nil, err
	}

	certASN, err := makeChildCert(parentCert, parentKey, myPriv.Public())
	if err != nil {
		return nil, err
	}

	attObj := protocol.AttestationObject{
		RawAuthData: authDataBytes,
		Format:      "fido-u2f",
		AttStatement: map[string]interface{}{
			"x5c": [][]byte{certASN},
			"sig": sig,
		},
	}

	attObjBytes, err := cbor.Marshal(attObj)
	if err != nil {
		return nil, err
	}

	// save credential data for future assertions
	credSource := pubKeyCredSource{
		ID:     base64.RawURLEncoding.EncodeToString(randCredID),
		Type:   "public-key",
		Priv:   myPriv,
		userID: params.Response.User.ID,
		RPID:   params.Response.RelyingParty.ID,
	}

	ma.CredSources[credSource.ID] = &credSource

	credResponse := credentialResponse{
		ID:   base64.RawURLEncoding.EncodeToString(authData.AttData.CredentialID),
		Type: "public-key",
		Response: responseData{
			AttestationObject: base64.RawURLEncoding.EncodeToString(attObjBytes),
			ClientDataJSON:    base64.RawURLEncoding.EncodeToString(clientDataJSON),
		},
	}

	return json.Marshal(credResponse)
}

func collectClientData(requesterOrigin string, challenge []byte) ([]byte, []byte, error) {
	collectedClientData := &protocol.CollectedClientData{
		Type:      "webauthn.create",
		Challenge: base64.RawURLEncoding.EncodeToString(challenge),
		Origin:    requesterOrigin,
	}

	clientDataJSON, err := json.Marshal(collectedClientData)
	if err != nil {
		return nil, nil, err
	}

	clientDataHash := sha256.Sum256(clientDataJSON)

	return clientDataJSON, clientDataHash[:], nil
}

// AssertionInternalResponse holds the internal data of an assertion response.
type AssertionInternalResponse struct {
	AuthenticatorData protocol.URLEncodedBase64 `json:"authenticatorData"`
	ClientDataJSON    protocol.URLEncodedBase64 `json:"clientDataJSON"`
	Signature         protocol.URLEncodedBase64 `json:"signature"`
}

// AssertionResponse holds an assertion response.
type AssertionResponse struct {
	ID       string                    `json:"id"`
	RawID    string                    `json:"rawId"`
	Type     string                    `json:"type"`
	Response AssertionInternalResponse `json:"response"`
}

// Assert creates a credential assertion for a relying party that's already registered via Authenticate.
// Returns a JSON object containing the assertion response.
func (ma *MockAuthenticator) Assert(requesterOrigin string, params *protocol.PublicKeyCredentialRequestOptions) ([]byte, error) { // nolint:lll // splitting lines would make this less readable
	collectedClientData := &protocol.CollectedClientData{
		Type:      "webauthn.get",
		Challenge: base64.RawURLEncoding.EncodeToString(params.Challenge),
		Origin:    requesterOrigin,
	}

	collectedClientDataBytes, err := json.Marshal(collectedClientData)
	if err != nil {
		return nil, err
	}

	ccdHash := sha256.Sum256(collectedClientDataBytes)

	credSource, err := ma.getValidCred(params.AllowedCredentials)
	if err != nil {
		return nil, err
	}

	rpIDHash := sha256.Sum256([]byte(params.RelyingPartyID))

	authBytes := marshalAuthData(&protocol.AuthenticatorData{
		RPIDHash: rpIDHash[:],
		Flags:    userPresentFlag | userVerifiedFlag, // bit 0: user present, bit 2: user verified
		Counter:  0,                                  // incrementing the counter is optional, so we won't
	}, false) // for credential assertion we don't include attestation data

	sigData := bytes.NewBuffer(nil)
	sigData.Write(authBytes)
	sigData.Write(ccdHash[:])

	sigHash := sha256.Sum256(sigData.Bytes())

	sig, err := credSource.Priv.Sign(rand.Reader, sigHash[:], nil)
	if err != nil {
		return nil, err
	}

	out := AssertionResponse{
		ID:    credSource.ID,
		RawID: credSource.ID,
		Type:  "public-key",
		Response: AssertionInternalResponse{
			AuthenticatorData: authBytes,
			ClientDataJSON:    collectedClientDataBytes,
			Signature:         sig,
		},
	}

	return json.Marshal(&out)
}

func (ma *MockAuthenticator) getValidCred(credOpts []protocol.CredentialDescriptor) (*pubKeyCredSource, error) {
	var credSource *pubKeyCredSource

	var ok bool

	for _, credOpt := range credOpts {
		credMapKey := base64.RawURLEncoding.EncodeToString(credOpt.CredentialID)
		if credSource, ok = ma.CredSources[credMapKey]; ok {
			break
		}
	}

	if !ok || credSource == nil || credSource.Priv == nil {
		return nil, fmt.Errorf("no usable cred sources")
	}

	return credSource, nil
}

// makeChildCert returns a certificate signed by parent, with its PEM encoding and private key.
func makeChildCert(parent *x509.Certificate, parentPriv, myPub interface{}) ([]byte, error) {
	certSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return nil, err
	}

	keyID := make([]byte, 16)
	if _, err = rand.Read(keyID); err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: certSerialNumber,
		Subject: pkix.Name{
			CommonName:   "BDD-test attestation certificate",
			SerialNumber: fmt.Sprint(*certSerialNumber),
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour * 24), // nolint:gomnd // really, nolint?
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		SubjectKeyId:          keyID,
	}

	return x509.CreateCertificate(rand.Reader, &template, parent, myPub, parentPriv)
}

// makeSig signs the authenticated data of an authentication statement.
func makeSig(rpIDHash, clientDataHash, credentialID []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	sigData := bytes.NewBuffer([]byte{0x00})
	sigData.Write(rpIDHash)
	sigData.Write(clientDataHash)
	sigData.Write(credentialID)
	sigData.Write([]byte{0x04}) // append octet 04, identifies uncompressed pubkey format
	sigData.Write(privKey.X.Bytes())
	sigData.Write(privKey.Y.Bytes())

	sigHash := sha256.Sum256(sigData.Bytes())

	return privKey.Sign(rand.Reader, sigHash[:], nil)
}

// marshalAuthData constructs the byte array of auth data in sections 6.1 and 6.5.1 of the spec.
//  https://www.w3.org/TR/webauthn/#sctn-authenticator-data
//  https://www.w3.org/TR/webauthn/#sctn-attested-credential-data
func marshalAuthData(authData *protocol.AuthenticatorData, includeAttData bool) []byte {
	authDataBytes := bytes.NewBuffer(nil)
	authDataBytes.Write(authData.RPIDHash)
	authDataBytes.Write([]byte{byte(authData.Flags)})

	counterBytes := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(counterBytes, authData.Counter)
	authDataBytes.Write(counterBytes)

	if includeAttData {
		authDataBytes.Write(authData.AttData.AAGUID)

		credIDLen := []byte{0, 0}
		binary.BigEndian.PutUint16(credIDLen, uint16(len(authData.AttData.CredentialID)))
		authDataBytes.Write(credIDLen)

		authDataBytes.Write(authData.AttData.CredentialID)
		authDataBytes.Write(authData.AttData.CredentialPublicKey)
	}

	return authDataBytes.Bytes()
}

func loadCert(filePath string) (*x509.Certificate, error) {
	data, errRead := ioutil.ReadFile(path.Clean(filePath))
	if errRead != nil {
		return nil, fmt.Errorf("failed to read cert file: %w", errRead)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode cert pem")
	}

	cert, errParse := x509.ParseCertificate(block.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse cert: %w", errParse)
	}

	return cert, nil
}

func loadKey(filePath string) (*ecdsa.PrivateKey, error) {
	data, errRead := ioutil.ReadFile(path.Clean(filePath))
	if errRead != nil {
		return nil, fmt.Errorf("failed to read key file: %w", errRead)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode key pem")
	}

	key, errParse := x509.ParseECPrivateKey(block.Bytes)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse key: %w", errParse)
	}

	return key, nil
}
