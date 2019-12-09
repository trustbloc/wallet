/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/trustbloc/edge-agent/cmd/http-server/startcmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "http server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	rootCmd.AddCommand(startcmd.GetStartCmd(&startcmd.HTTPServer{}))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run http server: %s", err.Error())
	}
}
