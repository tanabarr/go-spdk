// Go bindings for SPDK
package spdk

// #include <stdlib.h>
// #include "include/nvme_discover.h"
//
// struct ns_t*
// discover(void *f)
// {
// 	 struct ns_t* (*nvme_discover)();
//
// 	 nvme_discover = (struct ns_t* (*)())f;
// 	 return nvme_discover();
// }
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/coreos/pkg/dlopen"
	"github.com/pkg/errors"
)

// Namespace struct mirrors C.struct_ns_t and
// describes a NVMe Namespace tied to a controller.
//
// TODO: populate implicitly using inner member:
// +inner C.struct_ns_t
type Namespace struct {
	CtrlrModel  string
	CtrlrSerial string
	Id          int32
	Size        int32
}

// c2GoNamespace is a private translation function
func c2GoNamespace(ns *C.struct_ns_t) Namespace {
	return Namespace{
		CtrlrModel:  C.GoString(&ns.ctrlr_model[0]),
		CtrlrSerial: C.GoString(&ns.ctrlr_serial[0]),
		Id:          int32(ns.id),
		Size:        int32(ns.size),
	}
}

// rc2err returns an failure if rc != 0.
//
// TODO: If err is already set then it is wrapped,
// otherwise it is ignored. e.g.
// func rc2err(label string, rc C.int, err error) error {
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

// NVMeDiscover retrieves library handle and looks up
// symbol pointer for nvme_discover function.
// void point I cannot be called directly from Go
// so is passed to C function exposed by cgo which can
// cast as desired type and call.
// C.discover returns a pointer to single linked list
// of ns_t structs which are converted to a slice of
// Go Namespace structs.
func NVMeDiscover() []Namespace {
	var entries []Namespace

	lnvme_discover := []string{
		"spdk/libnvme_discover.so",
	}

	h, err := dlopen.GetHandle(lnvme_discover)
	if err != nil {
		fmt.Println(
			fmt.Errorf(
				`couldn't get a handle to the library %v: %v`,
				lnvme_discover,
				err))
		return
	}
	defer h.Close()
	fmt.Println(h)

	f := "nvme_discover"
	nvmeDiscover, err := h.GetSymbolPointer(f)
	if err != nil {
		fmt.Println(fmt.Errorf(`couldn't get symbol %q: %v`, f, err))
		return
	}
	fmt.Println(nvmeDiscover)

	ns_p := C.discover(nvmeDiscover)

	//ns_p := C.nvme_discover()
	//if err := rc2err("nvme_discover", rc); err != nil {
	//	return err
	//}

	for ns_p != nil {
		defer C.free(unsafe.Pointer(ns_p))
		entries = append(entries, c2GoNamespace(ns_p))
		ns_p = ns_p.next
	}

	return entries
}
