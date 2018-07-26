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

import "unsafe"

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

func main() {
	println("testing")
	entry_p := C.nvme_discover()
	defer C.free(unsafe.Pointer(entry_p))
	//retEntry := Entry{inner: *entry_p}
	retEntry := TranslateCEntry2GoEntry(entry_p)
	//	println(*retEntry.inner.name1)
	//	println(*retEntry.inner.name2)
	println(retEntry.name1)
	println(retEntry.name2)
}
