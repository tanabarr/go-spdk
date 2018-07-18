package spdk

/* Go bindings for SPDK */

// #cgo LDFLAGS: -L/usr/local/lib -lspdk
// #cgo CFLAGS: -I/usr/local/include
// #include "stdlib.h"
// #include "spdk/stdinc.h"
// #include "spdk/nvme.h"
// #include "spdk/env.h"
import "C"
import (
	"../stdio"
	//"fmt"
	"github.com/pkg/errors"
)

/**
 * \brief Enumerate the bus indicated by the transport ID and attach the userspace NVMe driver
 * to each device found if desired.
 *
 * \param trid The transport ID indicating which bus to enumerate. If the trtype is PCIe or trid is NULL,
 * this will scan the local PCIe bus. If the trtype is RDMA, the traddr and trsvcid must point at the
 * location of an NVMe-oF discovery service.
 * \param cb_ctx Opaque value which will be passed back in cb_ctx parameter of the callbacks.
 * \param probe_cb will be called once per NVMe device found in the system.
 * \param attach_cb will be called for devices for which probe_cb returned true once that NVMe
 * controller has been attached to the userspace driver.
 * \param remove_cb will be called for devices that were attached in a preyvious spdk_nvme_probe()
 * call but are no longer attached to the system. Optional; specify NULL if removal notices are not
 * desired.
 *
 * This function is not thread safe and should only be called from one thread at a time while no
 * other threads are actively using any NVMe devices.
 *
 * If called from a secondary process, only devices that have been attached to the userspace driver
 * in the primary process will be probed.
 *
 * If called more than once, only devices that are not already attached to the SPDK NVMe driver
 * will be reported.
 *
 * To stop using the the controller and release its associated resources,
 * call \ref spdk_nvme_detach with the spdk_nvme_ctrlr instance from the attach_cb() function.
 */
/**
 * Initialize the default value of opts.
 *
 * \param opts Data structure where SPDK will initialize the default options.
 */

// Returns an failure if rc != 0. If err is already set
// then it is wrapped, otherwise it is ignored.
//func rc2err(label string, rc C.int, err error) error {
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

func Init() error {
	stdio.Stdout.WriteString("Initializing NVMe Controllers\n")
	opts := &C.struct_spdk_env_opts{}
	//opts.name = "all"
	//stdio.Stdout.WriteString(C.GoString(opts.name) + "\n")
	//fmt.Printf("%v\n", *opts.name)
	rc := C.spdk_env_opts_init(opts)

	if err := rc2err("spdk_env_opts_init", rc); err != nil {
		return err
	}

	// stdio.defer C.free(unsafe.Pointer())
	// void spdk_env_opts_init(struct spdk_env_opts *opts);
	return nil
}

/**
 * Initialize the environment library. This must be called prior to using
 * any other functions in this library.
 *
 * \param opts Environment initialization options.
 * \return 0 on success, or negative errno on failure.
 */
// int spdk_env_init(const struct spdk_env_opts *opts);
//
// struct spdk_env_opts opts;

/*
 * SPDK relies on an abstraction around the local environment
 * named env that handles memory allocation and PCI device operations.
 * This library must be initialized first.
 *
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
//	var cGroup *C.char
//	if group != "" {
//		cGroup = C.CString(group)
//		defer C.free(unsafe.Pointer(cGroup))
//	}
//
//	cDev := C.CString("pmem")
//	defer C.free(unsafe.Pointer(cDev))
//
//	nranks := C.uint32_t(13)
//	svc := &C.daos_rank_list_t{}
//	ranks := allocRanks(nranks)
//	defer C.free(unsafe.Pointer(ranks))
//
//	svc.rl_nr.num = nranks
//	svc.rl_nr.num_out = 0
//	svc.rl_ranks = ranks
//
//	var u C.uuid_t
//	var uuid [unsafe.Sizeof(u)]C.uchar
//
//	rc, err := C.daos_pool_create(C.uint(mode),
//		C.uint(uid),
//		C.uint(gid),
//		cGroup,
//		nil, /* tgts */
//		cDev,
//		C.daos_size_t(size),
//		svc,
//		(*C.uchar)(unsafe.Pointer(&uuid[0])),
//		nil /* ev */)
//
//	if err = rc2err("daos_pool_create", rc, err); err != nil {
//		return "", err
//	}
//	return uuid2str(uuid[:]), nil
//}
//
//
