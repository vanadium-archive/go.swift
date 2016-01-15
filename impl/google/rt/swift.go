// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package rt

import (
	"unsafe"

	"v.io/v23"
	"v.io/v23/context"

	sutil "v.io/x/swift/util"
	scontext "v.io/x/swift/v23/context"
)

// #import "../../../types.h"
import "C"

type shutdownKey struct{}

//export swift_io_v_impl_google_rt_VRuntimeImpl_nativeInit
func swift_io_v_impl_google_rt_VRuntimeImpl_nativeInit() C.GoContextHandle {
	ctx, shutdownFunc := v23.Init()
	ctx = context.WithValue(ctx, shutdownKey{}, shutdownFunc)
	return C.GoContextHandle(scontext.SwiftContext(ctx))
}

//export swift_io_v_impl_google_rt_VRuntimeImpl_nativeShutdown
func swift_io_v_impl_google_rt_VRuntimeImpl_nativeShutdown(ctxHandle C.GoContextHandle) {
	ctx := scontext.GoContext(uint64(ctxHandle))
	value := ctx.Value(shutdownKey{})

	if shutdownFunc, ok := value.(v23.Shutdown); ok {
		shutdownFunc()
	}
}

//export swift_io_v_impl_google_rt_VRuntimeImpl_nativeWithNewClient
func swift_io_v_impl_google_rt_VRuntimeImpl_nativeWithNewClient(ctxHandle C.GoContextHandle, errOut *C.SwiftVError) C.GoContextHandle {
	ctx := scontext.GoContext(uint64(ctxHandle))
	// No options supported yet.
	newCtx, _, err := v23.WithNewClient(ctx)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoContextHandle(0)
	}

	return C.GoContextHandle(scontext.SwiftContext(newCtx))
}
