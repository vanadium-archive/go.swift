// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package rpc

import (
	"fmt"

	"v.io/v23/rpc"

	iutil "v.io/x/ios/util"
)

// #include "types.h"
import "C"

func SwiftClientCall(call rpc.ClientCall) C.GoClientCallHandle {
	return C.GoClientCallHandle(iutil.GoNewRef(call))
}

func GoClientCall(callHandle uint64) rpc.ClientCall {
	valptr := iutil.GoGetRef(callHandle)
	if call, ok := valptr.(rpc.ClientCall); ok {
		return call
	} else {
		panic(fmt.Sprintf("Couldn't get client call from handle with id %d", callHandle))
	}
}

func EmptySwiftByteArrayArray() C.SwiftByteArrayArray {
	var empty C.SwiftByteArrayArray
	empty.length = 0
	empty.data = nil
	return empty
}
