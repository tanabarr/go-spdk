// cmpout -tags=use_go_run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"../spdk"
	"fmt"
)

func main() {
	fmt.Println("Discovered NVMe devices: ")
	for _, e := range spdk.NVMeDiscover() {
		fmt.Printf(
			"controller (model/serial): %v/%v, namespace: %v, size: %v\n",
			e.CtrlrModel, e.CtrlrSerial, e.Id, e.Size)
	}
}
