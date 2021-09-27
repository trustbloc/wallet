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

	"github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/client/didexchange"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/presentproof"
	arieshttp "github.com/hyperledger/aries-framework-go/pkg/didcomm/transport/http"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/defaults"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/httpbinding"
)

type didComm struct {
	OOBClient          *outofband.Client
	DIDExchClient      *didexchange.Client
	PresentProofClient *presentproof.Client
}

func startAriesAgent() (*didComm, error) {
	var opts []aries.Option
	opts = append(opts, aries.WithStoreProvider(mem.NewProvider()))

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

	universalResolverVDRI, resErr := httpbinding.New(os.Getenv(resolverURL),
		httpbinding.WithAccept(acceptsDID), httpbinding.WithTLSConfig(tlsConfig))
	if resErr != nil {
		return nil, fmt.Errorf("failed to create new universal resolver vdri: %w", resErr)
	}

	opts = append(opts, aries.WithVDR(universalResolverVDRI))

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

	return &didComm{
		OOBClient:          oobClient,
		DIDExchClient:      didExClient,
		PresentProofClient: presentProofClient,
	}, nil
}

func acceptsDID(method string) bool {
	return method == "orb"
}
