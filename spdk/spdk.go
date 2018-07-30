// Go bindings for SPDK
package spdk

/** lib & include dirs detected from CGO_CFLAGS & CGO_LDFLAGS
 * env vars.
 */

/*
#cgo CFLAGS: -I /home/tanabarr/daos_m/_build.external/spdk/include -I .
#cgo LDFLAGS: -L /home/tanabarr/daos_m/_build.external/spdk/build/lib -L . -lnvme_discover -lspdk

#include "stdlib.h"
#include "spdk/stdinc.h"
#include "spdk/nvme.h"
#include "spdk/env.h"

#include "nvme_discover.h"
*/
import "C"

import (
	"github.com/pkg/errors"
	"unsafe"
	"fmt"
)

type NS struct {
	id int
	ctrlrName string
	// ctrlrSerial string
	size int
	//inner C.struct_ns_t
}

func TranslateCNS2GoNS(ns *C.struct_ns_t) NS {
	return NS{
		id: int(ns.id),
		ctrlrName: C.GoString(&ns.ctrlr_name[0]),
		size: int(ns.size),
	}
}

/** Returns an failure if rc != 0. If err is already set
 * then it is wrapped, otherwise it is ignored.
 *
 * //func rc2err(label string, rc C.int, err error) error {
 */
func rc2err(label string, rc C.int) error {
	if rc != 0 {
		if rc < 0 {
			rc = -rc
		}
		// e := errors.Error(rc)
		return errors.Errorf("%s: %s", label, rc) // e
	}
	return nil
}

/**
 * SPDK relies on an abstraction around the local environment
 * named env that handles memory allocation and PCI device operations.
 * This library must be initialized first.
 *
 * \return nil on success, err otherwise
 */
func InitSPDKEnv() error {
	println("Initializing NVMe Driver")
	opts := &C.struct_spdk_env_opts{}

	C.spdk_env_opts_init(opts)

	rc := C.spdk_env_init(opts)
	if err := rc2err("spdk_env_opts_init", rc); err != nil {
		return err
	}

	return nil
}

func NVMeDiscover() string {
	//devices_s := C.nvme_discover()
	//if err := rc2err("nvme_discover", rc); err != nil {
	//	return err
	//}
	var entries []NS
	ns_p := C.nvme_discover()

	for ns_p != nil {
		defer C.free(unsafe.Pointer(ns_p))
		entries = append(entries, TranslateCNS2GoNS(ns_p))
		ns_p = ns_p.next
	}
	for _, e := range entries {
		println(
			"controller: %v, namespace: %v, size: %v",
			entries[0].ctrlrName, entries[0].id, entries[0].size)
	}
	return fmt.Sprintf(
		"controller: %v, namespace: %v, size: %v\n",
		entries[0].ctrlrName, entries[0].id, entries[0].size)

	//return nil
}
