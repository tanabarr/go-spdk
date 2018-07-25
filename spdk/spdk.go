// Go bindings for SPDK
package spdk

/** lib & include dirs detected from CGO_CFLAGS & CGO_LDFLAGS
 * env vars.
 */

/* 
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lnvme_discover -lspdk

#include "stdlib.h"
#include "spdk/stdinc.h"
#include "spdk/nvme.h"
#include "spdk/env.h"

#include "nvme_discover.h"
*/
import "C"

import (
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

    C.spdk_env_opts_init(opts)
    
    rc := C.spdk_env_init(opts)
    if err := rc2err("spdk_env_opts_init", rc); err != nil {
    	return err
    }
    
    return nil
}

func NVMeDiscover() error {
    rc := C.nvme_discover()
    if err := rc2err("nvme_discover", rc); err != nil {
    	return err
    }

    return nil
}
