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
)

type NameSpace struct {
	Id int32
	CtrlrName string
	// ctrlrSerial string
	Size int32
	//inner C.struct_ns_t
}

func c2GoNameSpace(ns *C.struct_ns_t) NameSpace {
	return NameSpace{
		Id: int32(ns.id),
		CtrlrName: C.GoString(&ns.ctrlr_name[0]),
		Size: int32(ns.size),
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

func NVMeDiscover() []NameSpace {
	var entries []NameSpace
	ns_p := C.nvme_discover()
	//if err := rc2err("nvme_discover", rc); err != nil {
	//	return err
	//}

	for ns_p != nil {
		defer C.free(unsafe.Pointer(ns_p))
		entries = append(entries, c2GoNameSpace(ns_p))
		ns_p = ns_p.next
	}
	return entries
}
