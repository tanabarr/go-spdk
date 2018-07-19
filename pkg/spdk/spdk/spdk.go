// Go bindings for SPDK
package spdk

// #cgo LDFLAGS: -L/usr/local/lib -lspdk
// #cgo CFLAGS: -I/usr/local/include
// #include "stdlib.h"
// #include "spdk/stdinc.h"
// #include "spdk/nvme.h"
// #include "spdk/env.h"
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

/**
	println("Initializing NVMe Controllers")
	// stdio.defer C.free(unsafe.Pointer())
	// void spdk_env_opts_init(struct spdk_env_opts *opts);
 */

// spdk_env_opts_init(&opts);
// opts.name = "hello_world";
// opts.shm_id = 0;
// if (spdk_env_init(&opts) < 0) {
// fprintf(stderr, "Unable to initialize SPDK env\n");
// return 1;
// }
//printf("Initializing NVMe Controllers\n");

// func NVMEProbe(mode uint32, uid uint32, gid uint32, group string, size int64) (string, error) {
//    C.spdk_nvme_probe()
//
////int spdk_nvme_probe(const struct spdk_nvme_transport_id *trid,
////            void *cb_ctx,
////            spdk_nvme_probe_cb probe_cb,
////            spdk_nvme_attach_cb attach_cb,
////            spdk_nvme_remove_cb remove_cb);