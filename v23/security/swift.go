// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

// Package security implements the security API surrounding principals and blessings
package security

import (
	"encoding/base64"
	"fmt"
	"unsafe"

	"v.io/v23"
	"v.io/v23/security"
	"v.io/v23/vom"
	"v.io/x/swift/util"
	"v.io/x/swift/v23/context"
)

/*
#include <string.h> // memcpy
#import "../../types.h"
*/
import "C"

//export swift_io_v_v23_security_simple_nativeGetPublicKey
func swift_io_v_v23_security_simple_nativeGetPublicKey(ctxHandle C.GoContextHandle, errOut *C.SwiftVError) *C.char {
	ctx := context.GoContext(uint64(ctxHandle))
	der, err := v23.GetPrincipal(ctx).PublicKey().MarshalBinary()
	if err != nil {
		util.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return nil
	}
	// Swift will have to free the allocation
	return C.CString(base64.URLEncoding.EncodeToString(der))
}

//export swift_io_v_v23_security_simple_nativeSetDefaultBlessings
func swift_io_v_v23_security_simple_nativeSetDefaultBlessings(ctxHandle C.GoContextHandle, encodedSwiftBlessings C.SwiftByteArray, errOut *C.SwiftVError) {
	ctx := context.GoContext(uint64(ctxHandle))
	encodedBlessings := util.GoBytesNoCopy(unsafe.Pointer(&encodedSwiftBlessings))
	var blessings security.Blessings
	if err := vom.Decode(encodedBlessings, &blessings); err != nil {
		ctx.Error("Unable to decode:", err)
		util.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
	principal := v23.GetPrincipal(ctx)
	if err := principal.BlessingStore().SetDefault(blessings); err != nil {
		util.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
}

//export swift_io_v_v23_security_simple_nativeGetDefaultBlessingsDebugString
func swift_io_v_v23_security_simple_nativeGetDefaultBlessingsDebugString(ctxHandle C.GoContextHandle) *C.char {
	ctx := context.GoContext(uint64(ctxHandle))
	blessings, _ := v23.GetPrincipal(ctx).BlessingStore().Default()
	return C.CString(fmt.Sprintf("%v", blessings))
}
