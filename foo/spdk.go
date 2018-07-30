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
	"unsafe"
)

type Entry struct {
	name1 string
	name2 string
	//inner C.struct_entry_t
}

func TranslateCEntry2GoEntry(e *C.struct_entry_t) Entry {
	return Entry{
		name1: C.GoString(&e.name1[0]),
		name2: C.GoString(&e.name2[0]),
	}
}

func main() {
	println("testing")
	var entries []Entry
	//devMap := make(map[string][]Entry)
	entry_p := C.nvme_discover()

	for entry_p != nil {
		defer C.free(unsafe.Pointer(entry_p))
		entries = append(entries, TranslateCEntry2GoEntry(entry_p))
		// devMap[retEntry.name1] = retEntry
		entry_p = entry_p.next
	}
	for _, e := range entries {
		println(e.name2)
		//"len=%d cap=%d %v\n", len(entries), cap(entries), entries)
	}
}

