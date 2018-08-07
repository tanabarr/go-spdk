// cmpout -tags=use_go_run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package main

// #include <string.h>
// #include <stdlib.h>
//
// int
// my_init(void *seoi, void *sei, void *sseo)
// {
//	//struct (*sseo) opts;
//	/*
//	 * SPDK relies on an abstraction around the local environment
//	 * named env that handles memory allocation and PCI device operations.
//	 * This library must be initialized first.
//	 *
//	 */
//	//(*seoi)(&opts);
//	//opts.name = "hello_world";
//	//opts.shm_id = 0;
//	//if ((*sei)(&opts) < 0) {
//	//	fprintf(stderr, "Unable to initialize SPDK env\n");
//	//	return 1;
//	//}
//	return 0;
// }
import "C"

import (
	"os"
	"fmt"
	//"unsafe"
	"github.com/coreos/pkg/dlopen"
)

//import "fmt"
//import "os"
//import "../spdk"
//import "github.com/coreos/pkg/dlopen"


func main() {
	// determine plugin to load
	//mod := "/home/tanabarr/daos_m/install/lib/libspdk.so"
	
	lspdk := []string{
		"/home/tanabarr/daos_m/install/lib/libspdk.so",
	}
	lnvme_discover := []string{
		"spdk/libnvme_discover.so",
//		"libc.so.6",
//		"libc.so",
	}
	fmt.Println(lspdk)
	fmt.Println(lnvme_discover)

	//h, err := dlopen.GetHandle(lspdk)
	h, err := dlopen.GetHandle(lnvme_discover)
	if err != nil {
		fmt.Println(fmt.Errorf(`couldn't get a handle to the library: %v`, err))
		os.Exit(1)
	}
	defer h.Close()
	fmt.Println(h)

	f := "nvme_discover"
	nvmeDiscover, err := h.GetSymbolPointer(f)
	if err != nil {
		fmt.Println(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
		os.Exit(1)
	}
	fmt.Println(nvmeDiscover)

	//f := "spdk_env_opts_init"
	//spdkOptsInit, err := h.GetSymbolPointer(f)
	//if err != nil {
	//	fmt.Println(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
	//	os.Exit(1)
	//}
	//fmt.Println(spdkOptsInit)

	//f = "spdk_env_opts"
	//spdkOpts, err := h.GetSymbolPointer(f)
	//if err != nil {
	//	fmt.Println(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
	//	os.Exit(1)
	//}
	//fmt.Println(spdkOpts)

	//f = "spdk_env_init"
	//spdkInit, err := h.GetSymbolPointer(f)
	//if err != nil {
	//	fmt.Println(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
	//	os.Exit(1)
	//}
	//fmt.Println(spdkInit)

	//opts := &C.spdk_env_opts_init{} //spdkOptsInit{}
	rc := C.my_init(nil, nil, nil)
	//rc := C.my_init(spdkOptsInit, spdkInit, spdkOpts)
	fmt.Println(rc)

	//return int(len), nil
	//fmt.Println("length of string is %d", len)

	os.Exit(0)
}

//	// load module
//	// 1. open the so file to load the symbols
//	plug, err := plugin.Open(mod)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	// 2. look up a symbol (an exported function or variable)
//	// in this case, variable Greeter
//	spdkEnvOptsInit, err := plug.Lookup("spdk_env_opts_init")
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	fmt.Println(spdkEnvOptsInit)
//	// 3. Assert that loaded symbol is of a desired type
//	// in this case interface type Greeter (defined above)
//	//var greeter Greeter
//	//greeter, ok := symGreeter.(Greeter)
//	//if !ok {
//	//	fmt.Println("unexpected type from module symbol")
//	//	os.Exit(1)
//	//}
//
//	// 4. use the module
//	//greeter.Greet()

//}
//func main() {
//	if err := spdk.InitSPDKEnv(); err != nil {
//		fmt.Printf("Unable to initialise SPDK env (%s)\n", err)
//	}
//
//	fmt.Println("Discovered NVMe devices: ")
//	for _, e := range spdk.NVMeDiscover() {
//		fmt.Printf(
//			"controller (model/serial): %v/%v, namespace: %v, size: %v\n",
//			e.CtrlrModel, e.CtrlrSerial, e.Id, e.Size)
//	}
//}
