/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("user wasm start")

	done := make(chan struct{})

	// TODO clean up custom wasm code

	<-done
}
