// Go bindings for SPDK
package main

/** lib & include dirs detected from CGO_CFLAGS & CGO_LDFLAGS
 * env vars.
 */

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lnvme_discover

#include "stdlib.h"

#include "nvme_discover.h"
*/
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Entry struct {
	name1 string
	name2 string
	//inner C.struct_entry_t
}

func TranslateCEntry2GoEntry(e *C.struct_entry_t) *Entry {
	if e == nil {
		return nil
	}
	return &Entry{
		name1: C.GoString(e.name1),
		name2: C.GoString(e.name2),
	}
}

func LogIfErr(err error) error {
	if err != nil {
		errStatus, _ := status.FromError(err)
		function, file, line, _ := runtime.Caller(1)

		// replace with new elaborated error
		err = status.Errorf(
			errStatus.Code(),
			fmt.Sprintf(
				"%v:l%d - %v(_), %v",
				file, line,
				runtime.FuncForPC(function).Name(),
				errStatus.Message()))

		log.Println(err)
	}
	return err
}

func main() {
	println("testing")
	entry_p := C.nvme_discover()
	defer C.free(unsafe.Pointer(entry_p))
	//retEntry := Entry{inner: *entry_p}
	retEntry := TranslateCEntry2GoEntry(entry_p)
	//	println(*retEntry.inner.name1)
	//	println(*retEntry.inner.name2)
	LogIfErr(status.Errorf(codes.InvalidArgument, "something went wrong"))
	LogIfErr(nil)
	//fmt.Println("%s %s %s", runtime.Caller(1))
	println(retEntry.name1)
	println(retEntry.name2)
}
