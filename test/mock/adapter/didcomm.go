/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto"
	"crypto/ed25519"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/hyperledger/aries-framework-go-ext/component/vdr/orb"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/issuecredential"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofbandv2"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	"github.com/hyperledger/aries-framework-go/pkg/didcomm/transport"
	arieshttp "github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/http"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk/jwksupport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext/remote"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/jwkkid"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/defaults"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	tlsutils "github.com/trustbloc/edge-core/pkg/utils/tls"
)

type didComm struct {
	OOBClient             *outofband.Client
	OOBV2Client           *outofbandv2.Client
	DIDExchClient         *didexchange.Client
	PresentProofClient    *presentproof.Client
	IssueCredentialClient *issuecredential.Client
	OrbDIDV2              string
}

var (
	//nolint:gochecknoglobals // translation tables copied from afgo for consistency
	keyTypes = map[string]kms.KeyType{
		"ed25519":           kms.ED25519Type,
		"ecdsap256ieee1363": kms.ECDSAP256TypeIEEEP1363,
		"ecdsap256der":      kms.ECDSAP256TypeDER,
		"ecdsap384ieee1363": kms.ECDSAP384TypeIEEEP1363,
		"ecdsap384der":      kms.ECDSAP384TypeDER,
		"ecdsap521ieee1363": kms.ECDSAP521TypeIEEEP1363,
		"ecdsap521der":      kms.ECDSAP521TypeDER,
	}

	//nolint:gochecknoglobals // translation tables copied from afgo for consistency
	keyAgreementTypes = map[string]kms.KeyType{
		"x25519kw": kms.X25519ECDHKWType,
		"p256kw":   kms.NISTP256ECDHKWType,
		"p384kw":   kms.NISTP384ECDHKWType,
		"p521kw":   kms.NISTP521ECDHKWType,
	}
)

func startAriesAgent() (*didComm, error) {
	storeProvider := mem.NewProvider()

	var opts []aries.Option
	opts = append(opts, aries.WithStoreProvider(storeProvider))

	opts = append(opts, defaults.WithInboundHTTPAddr(os.Getenv(didCommInternalHostEnvKey),
		os.Getenv(didCommExternalHostEnvKey), os.Getenv(tlsCertFileEnvKey),
		os.Getenv(tlsKeyFileEnvKey)))

	tlsConfig := &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12}

	outbound, err := arieshttp.NewOutbound(arieshttp.WithOutboundHTTPClient(
		&http.Client{Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		}}),
	)
	if err != nil {
		return nil, fmt.Errorf("http outbound transport initialization failed: %w", err)
	}

	opts = append(opts, aries.WithOutboundTransports(outbound))

	vdri, err := orb.New(nil,
		orb.WithTLSConfig(tlsConfig),
		orb.WithDomain(os.Getenv(orbDomainEnvKey)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init orb VDR: %w", err)
	}

	var (
		useDIDCommV2 bool
		keyT         kms.KeyType
		keyAgrT      kms.KeyType
	)

	useDIDCommV2, err = strconv.ParseBool(os.Getenv(useDIDCommV2EnvKey))
	if useDIDCommV2 && err == nil {
		opts = append(opts, aries.WithMediaTypeProfiles([]string{transport.MediaTypeDIDCommV2Profile}))
	}

	// if keyTypeEnvKey and KeyAgreementTypeEnvKey are found, then override agent's options
	if kt, ok := os.LookupEnv(keyTypeEnvKey); ok {
		keyT = keyTypes[kt]
		opts = append(opts, aries.WithKeyType(keyT))
	}

	if kt, ok := os.LookupEnv(keyAgreementTypeEnvKey); ok {
		keyAgrT = keyAgreementTypes[kt]
		opts = append(opts, aries.WithKeyAgreementType(keyAgrT))
	}

	opts = append(opts, aries.WithVDR(vdri))

	if ctxURL := os.Getenv(contextProviderEnvKey); ctxURL != "" {
		if caCerts := os.Getenv(tlsCACertsEnvKey); caCerts != "" {
			rootCAs, err := tlsutils.GetCertPool(true, []string{caCerts})
			if err != nil {
				panic("failed to setup root ca, " + err.Error())
			}

			docLoader, err := createJSONLDDocumentLoader(storeProvider, &tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12}, ctxURL)
			if err != nil {
				panic("failed to setup document loader, " + err.Error())
			}

			opts = append(opts, aries.WithJSONLDDocumentLoader(docLoader))
		} else {
			opts = append(opts, aries.WithJSONLDContextProviderURL(ctxURL))
		}
	}

	framework, err := aries.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize framework :  %w", err)
	}

	ctx, err := framework.Context()
	if err != nil {
		return nil, fmt.Errorf("failed to get aries context : %w", err)
	}

	// out-of-band client
	oobClient, err := outofband.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create oob-client : %w", err)
	}

	publicDIDV2 := ""

	if useDIDCommV2 {
		publicDIDV2, err = createPublicDIDV2(vdri, ctx.KMS(), keyT, keyAgrT)
		if err != nil {
			return nil, fmt.Errorf("failed to create orb DID for OOB V2 invitations: %w", err)
		}
	}

	// out-of-band v2 client
	oobV2Client, err := outofbandv2.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create oob-client V2 : %w", err)
	}

	// did-exchange client
	didExClient, err := didexchange.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create didexchange-client : %w", err)
	}

	// present-proof client
	presentProofClient, err := presentproof.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create presentproof-client: %w", err)
	}

	// issue-credential client
	issueCredentialClient, err := issuecredential.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create issuecredential-client: %w", err)
	}

	return &didComm{
		OOBClient:             oobClient,
		OOBV2Client:           oobV2Client,
		DIDExchClient:         didExClient,
		PresentProofClient:    presentProofClient,
		IssueCredentialClient: issueCredentialClient,
		OrbDIDV2:              publicDIDV2,
	}, nil
}

func createPublicDIDV2(vdri vdr.VDR, km kms.KeyManager, keyType, keyAgreementType kms.KeyType) (string, error) {
	didDoc, err := buildDIDDocV2(km, keyType, keyAgreementType)
	if err != nil {
		return "", fmt.Errorf("failed to create DID doc: %w", err)
	}

	updateKey, err := newKey(km)
	if err != nil {
		return "", fmt.Errorf("failed to create udpateKey for vdri.Create(): %w", err)
	}

	recoveryKey, err := newKey(km)
	if err != nil {
		return "", fmt.Errorf("failed to create recoveryKey for vdri.Create(): %w", err)
	}

	docRes, err := vdri.Create(didDoc, vdr.WithOption(orb.UpdatePublicKeyOpt, updateKey),
		vdr.WithOption(orb.RecoveryPublicKeyOpt, recoveryKey))
	if err != nil {
		return "", fmt.Errorf("failed to create orb DID from VDRI: %w", err)
	}

	return docRes.DIDDocument.ID, nil
}

func newKey(km kms.KeyManager) (crypto.PublicKey, error) {
	_, bits, err := km.CreateAndExportPubKeyBytes(kms.ED25519Type)
	if err != nil {
		return nil, fmt.Errorf("failed to create key : %w", err)
	}

	return ed25519.PublicKey(bits), nil
}

func buildDIDDocV2(km kms.KeyManager, keyType, keyAgreementType kms.KeyType) (*did.Doc, error) {
	didDoc := did.Doc{}

	auth, err := createVerification("#key-1", km, keyType, did.Authentication)
	if err != nil {
		return nil, fmt.Errorf("creating did doc Authentication: %w", err)
	}

	didDoc.Authentication = append(didDoc.Authentication, *auth)

	kagr, err := createVerification("#key-2", km, keyAgreementType, did.KeyAgreement)
	if err != nil {
		return nil, fmt.Errorf("creating did doc KeyAgreement: %w", err)
	}

	didDoc.KeyAgreement = append(didDoc.KeyAgreement, *kagr)

	didDoc.Service = []did.Service{{
		ID:              uuid.NewString(),
		ServiceEndpoint: os.Getenv(didCommExternalHostEnvKey),
		Type:            "DIDCommMessaging",
	}}

	return &didDoc, nil
}

func createVerification(id string, km kms.KeyManager, kt kms.KeyType, relationship did.VerificationRelationship,
) (*did.Verification, error) {
	vm, err := createVerificationMethod(id, km, kt)
	if err != nil {
		return nil, fmt.Errorf("creating verification: %w", err)
	}

	return did.NewReferencedVerification(vm, relationship), nil
}

func createVerificationMethod(id string, km kms.KeyManager, kt kms.KeyType) (*did.VerificationMethod, error) {
	kid, pkBytes, err := km.CreateAndExportPubKeyBytes(kt)
	if err != nil {
		return nil, fmt.Errorf("creating public key: %w", err)
	}

	var j *jwk.JWK

	if kt == kms.ED25519Type {
		j, err = jwksupport.JWKFromKey(ed25519.PublicKey(pkBytes))
		if err != nil {
			return nil, fmt.Errorf("converting ed25519 key to JWK: %w", err)
		}

		id = kid
	} else {
		j, err = jwkkid.BuildJWK(pkBytes, kt)
		if err != nil {
			return nil, fmt.Errorf("creating JWK: %w", err)
		}

		j.KeyID = kid
	}

	jwk2020 := "JsonWebKey2020"

	vm, err := did.NewVerificationMethodFromJWK(id, jwk2020, "", j)
	if err != nil {
		return nil, fmt.Errorf("creating verification method: %w", err)
	}

	return vm, nil
}

func createJSONLDDocumentLoader(store storage.Provider, tlsConfig *tls.Config,
	providerURL string) (*ld.DocumentLoader, error) {
	var loaderOpts []ld.DocumentLoaderOpts

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	if providerURL != "" {
		loaderOpts = append(loaderOpts,
			ld.WithRemoteProvider(
				remote.NewProvider(providerURL, remote.WithHTTPClient(httpClient)),
			),
		)
	}

	ldStore, err := NewLDStoreProvider(store)
	if err != nil {
		return nil, fmt.Errorf("failed to init LD store provider: %w", err)
	}

	loader, err := ld.NewDocumentLoader(ldStore, loaderOpts...)
	if err != nil {
		return nil, fmt.Errorf("new document loader: %w", err)
	}

	return loader, nil
}
