/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/spf13/cobra"
	"github.com/trustbloc/edge-agent/cmd/wallet-server/startcmd"
	"github.com/trustbloc/edge-core/pkg/log"
)

var logger = log.New("wallet-server")

func main() {
	rootCmd := &cobra.Command{
		Use: "wallet-server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	rootCmd.AddCommand(startcmd.GetStartCmd(&startcmd.HTTPServer{}))

	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Failed to run http server: %s", err.Error())
	}
}
