// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package rt

import (
	"unsafe"

	"v.io/v23"
	"v.io/v23/context"

	iutil "v.io/x/ios/util"
	icontext "v.io/x/ios/v23/context"
)

// #import "types.h"
import "C"

type shutdownKey struct{}

//export ios_io_v_impl_google_rt_VRuntimeImpl_nativeInit
func ios_io_v_impl_google_rt_VRuntimeImpl_nativeInit() C.GoContextHandle {
	ctx, shutdownFunc := v23.Init()
	ctx = context.WithValue(ctx, shutdownKey{}, shutdownFunc)
	return C.GoContextHandle(icontext.SwiftContext(ctx))
}

//export ios_io_v_impl_google_rt_VRuntimeImpl_nativeShutdown
func ios_io_v_impl_google_rt_VRuntimeImpl_nativeShutdown(ctxHandle C.GoContextHandle) {
	ctx := icontext.GoContext(uint64(ctxHandle))
	value := ctx.Value(shutdownKey{})

	if shutdownFunc, ok := value.(v23.Shutdown); ok {
		shutdownFunc()
	}
}

//export ios_io_v_impl_google_rt_VRuntimeImpl_nativeWithNewClient
func ios_io_v_impl_google_rt_VRuntimeImpl_nativeWithNewClient(ctxHandle C.GoContextHandle, errOut *C.SwiftVError) C.GoContextHandle {
	ctx := icontext.GoContext(uint64(ctxHandle))
	// No options supported yet.
	newCtx, _, err := v23.WithNewClient(ctx)
	if err != nil {
		iutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoContextHandle(0)
	}

	return C.GoContextHandle(icontext.SwiftContext(newCtx))
}
