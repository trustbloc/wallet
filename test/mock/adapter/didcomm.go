/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/hyperledger/aries-framework-go/pkg/client/issuecredential"

	"github.com/hyperledger/aries-framework-go-ext/component/vdr/orb"
	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	arieshttp "github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/http"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext/remote"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/defaults"
	"github.com/hyperledger/aries-framework-go/spi/storage"
	tlsutils "github.com/trustbloc/edge-core/pkg/utils/tls"
)

type didComm struct {
	OOBClient             *outofband.Client
	DIDExchClient         *didexchange.Client
	PresentProofClient    *presentproof.Client
	IssueCredentialClient *issuecredential.Client
}

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
		DIDExchClient:         didExClient,
		PresentProofClient:    presentProofClient,
		IssueCredentialClient: issueCredentialClient,
	}, nil
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
