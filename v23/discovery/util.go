// Copyright 2016 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package discovery

import (
	"fmt"
	"unsafe"

	"v.io/v23/discovery"
	"v.io/x/swift/util"
)

/*
#include <string.h> // memcpy
#include <stdlib.h>
#import "../../types.h"
*/
import "C"

// Converts a native pointer from Swift to a Go Discovery.T. Panics if it can't.
func GoDiscoveryT(discoveryHandle uint64) discovery.T {
	valptr := util.GoGetRef(discoveryHandle)
	if d, ok := valptr.(*discovery.T); ok {
		return *d
	} else {
		panic(fmt.Sprintf("Couldn't get discovery.T from handle with id %d", discoveryHandle))
	}
}

func swiftBytesCopy(data []byte) C.SwiftByteArray {
	var a C.SwiftByteArray
	a.length = C._GoUint64(len(data))
	a.data = C.malloc(C.size_t(len(data)))
	if a.data == nil {
		panic(fmt.Errorf("Unable to allocate %d bytes", a.length))
	}
	C.memmove(a.data, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	return a
}

func emptySwiftByteArray() C.SwiftByteArray {
	var empty C.SwiftByteArray
	empty.length = 0
	empty.data = nil
	return empty
}

func (s *C.SwiftCString) toString() string {
	return C.GoStringN(s.data, C.int(s.length))
}

func (a *C.SwiftCStringArray) toStrings() []string {
	var cstrs = (*[1 << 30]C.SwiftCString)(unsafe.Pointer(a.data))[:int(a.length):int(a.length)]
	var ret []string
	for _, cstr := range cstrs {
		ret = append(ret, cstr.toString())
	}
	return ret
}
