// cmpout -tags=use_go_run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build test_run

package main

import "fmt"
import "../spdk"

func main() {
	if err := spdk.InitSPDKEnv(); err != nil {
		fmt.Printf("Unable to initialise SPDK env (%s)\n", err)
	}

	fmt.Println("Discovered NVMe devices: ")
	for _, e := range spdk.NVMeDiscover() {
		fmt.Printf(
			"controller: %v, namespace: %v, size: %v\n",
			e.ctrlrName, e.id, e.size)
	}
}
