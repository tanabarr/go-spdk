// Go bindings for SPDK
package spdk

/** lib & include dirs detected from CGO_CFLAGS & CGO_LDFLAGS
 * env vars.
 */

/* 
#cgo LDFLAGS: -lspdk
#include "stdlib.h"
#include "spdk/stdinc.h"
#include "spdk/nvme.h"
#include "spdk/env.h"

typedef void (*probe_func)();
void nvme_probe(probe_func f) {
	f(); // NULL, NULL, NULL, NULL, NULL);
}
*/
import "C"

import (
	"unsafe"
	"github.com/pkg/errors"
)

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
	//opts.name = "all"
	//stdio.Stdout.WriteString(C.GoString(opts.name) + "\n")
	//fmt.Printf("%v\n", *opts.name)
	C.spdk_env_opts_init(opts)

	rc := C.spdk_env_init(opts)
	if err := rc2err("spdk_env_opts_init", rc); err != nil {
		return err
	}

	return nil
}

func NVMeProbe() {
    println("Initializing NVMe Controllers")
	C.nvme_probe((C.probe_func)(unsafe.Pointer(C.spdk_nvme_probe)))
}

////int spdk_nvme_probe(const struct spdk_nvme_transport_id *trid,
////            void *cb_ctx,
////            spdk_nvme_probe_cb probe_cb,
////            spdk_nvme_attach_cb attach_cb,
////            spdk_nvme_remove_cb remove_cb);
