//
// (C) Copyright 2018 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

// Package nvme provides Go bindings for SPDK: NVMe tasks
package nvme

// CGO_CFLAGS & CGO_LDFLAGS env vars can be used
// to specify additional dirs.

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lnvme_discover -lspdk

#include "stdlib.h"
#include "spdk/stdinc.h"
#include "spdk/nvme.h"
#include "spdk/env.h"

#include "include/nvme_discover.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// Controller struct mirrors C.struct_ctrlr_t and
// describes a NVMe controller.
//
// TODO: populate implicitly using inner member:
// +inner C.struct_ctrlr_t
type Controller struct {
	ID      int32
	Model   string
	Serial  string
	PCIAddr string
}

// Namespace struct mirrors C.struct_ns_t and
// describes a NVMe Namespace tied to a controller.
//
// TODO: populate implicitly using inner member:
// +inner C.struct_ns_t
type Namespace struct {
	ID    int32
	Size  int32
	Ctrlr *Controller
}

// c2GoController is a private translation function
func c2GoController(ctrlr *C.struct_ctrlr_t) Controller {
	return Controller{
		ID:      int32(ctrlr.id),
		Model:   C.GoString(&ctrlr.model[0]),
		Serial:  C.GoString(&ctrlr.serial[0]),
		PCIAddr: C.GoString(&ctrlr.serial[0]),
	}
}

// c2GoNamespace is a private translation function
func c2GoNamespace(ns *C.struct_ns_t, ctrlr *Controller) Namespace {
	return Namespace{
		ID:    int32(ns.id),
		Size:  int32(ns.size),
		Ctrlr: ctrlr,
	}
}

// rc2err returns an failure if rc != 0.
//
// TODO: If err is already set then it is wrapped,
// otherwise it is ignored. e.g.
// func rc2err(label string, rc C.int, err error) error {
//
// \return nil on success, err otherwise
func rc2err(label string, rc C.int) error {
	if rc != 0 {
		if rc < 0 {
			rc = -rc
		}
		// e := errors.Error(rc)
		return fmt.Errorf("%s: %s", label, rc) // e
	}
	return nil
}

// InitSPDKEnv initializes the SPDK environment.
//
// SPDK relies on an abstraction around the local environment
// named env that handles memory allocation and PCI device operations.
// This library must be initialized first.
//
// \return nil on success, err otherwise
func InitSPDKEnv() error {
	opts := &C.struct_spdk_env_opts{}

	C.spdk_env_opts_init(opts)

	rc := C.spdk_env_init(opts)
	if err := rc2err("spdk_env_opts_init", rc); err != nil {
		return err
	}

	return nil
}

// Discover calls C.nvme_discover which returns
// pointers to single linked list of ctrlr_t and ns_t structs.
// These are converted to slices of Controller and Namespace structs.
//
// \return ([]Controllers, []Namespace, nil) on success,
//         (nil, nil, error) otherwise
func Discover() ([]Controller, []Namespace, error) {
	if retPtr := C.nvme_discover(); retPtr != nil {
		var ctrlrs []Controller
		var nss []Namespace

		if retPtr.success == true {
			ctrlrPtr := retPtr.ctrlrs
			for ctrlrPtr != nil {
				defer C.free(unsafe.Pointer(ctrlrPtr))
				ctrlrs = append(ctrlrs, c2GoController(ctrlrPtr))
				ctrlrPtr = ctrlrPtr.next
			}

			return ctrlrs, nss, nil
		}

		return nil, nil, fmt.Errorf(
			"NVMeDiscover(): C.nvme_discover failed, verify SPDK install")
	}

	return nil, nil, fmt.Errorf(
		"NVMeDiscover(): C.nvme_discover unexpectedly returned NULL")
}
