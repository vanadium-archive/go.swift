// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package context

import (
	"unsafe"

	"v.io/v23/context"
	sutil "v.io/x/swift/util"
)

//#import "../../types.h"
import "C"

// Not currently used
//export swift_io_v_v23_context_VContext_nativeDeadline
func swift_io_v_v23_context_VContext_nativeDeadline(ctxHandle C.GoContextHandle) C.double {
	ctx := GoContext(uint64(ctxHandle))
	d, ok := ctx.Deadline()
	if !ok {
		return 0
	}
	return C.double(sutil.NSTimeInterval(d))
}

//export swift_io_v_v23_context_VContext_nativeWithCancel
func swift_io_v_v23_context_VContext_nativeWithCancel(ctxHandle C.GoContextHandle, errOut C.SwiftVErrorPtr) C.GoCancelableContextHandle {
	ctx := GoContext(uint64(ctxHandle))
	ctx, cancelFunc := context.WithCancel(ctx)
	swiftCtx, err := SwiftCancelableContext(ctx, cancelFunc)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoCancelableContextHandle(0)
	}
	return C.GoCancelableContextHandle(swiftCtx)
}

//export swift_io_v_v23_context_VContext_nativeWithDeadline
func swift_io_v_v23_context_VContext_nativeWithDeadline(ctxHandle C.GoContextHandle, deadlineEpoch C.double, errOut C.SwiftVErrorPtr) C.GoCancelableContextHandle {
	ctx := GoContext(uint64(ctxHandle))
	deadline := sutil.GoTime(float64(deadlineEpoch))
	ctx, cancelFunc := context.WithDeadline(ctx, deadline)
	swiftCtx, err := SwiftCancelableContext(ctx, cancelFunc)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoCancelableContextHandle(0)
	}
	return C.GoCancelableContextHandle(swiftCtx)
}

//export swift_io_v_v23_context_VContext_nativeWithTimeout
func swift_io_v_v23_context_VContext_nativeWithTimeout(ctxHandle C.GoContextHandle, nsTimeout C.double, errOut *C.SwiftVError) C.GoCancelableContextHandle {
	ctx := GoContext(uint64(ctxHandle))
	timeout := sutil.GoDuration(float64(nsTimeout))
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	swiftCtx, err := SwiftCancelableContext(ctx, cancelFunc)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoCancelableContextHandle(0)
	}
	return C.GoCancelableContextHandle(swiftCtx)
}

//export swift_io_v_v23_context_VContext_nativeFinalize
func swift_io_v_v23_context_VContext_nativeFinalize(ctxHandle C.GoContextHandle) {
	sutil.GoUnref(uint64(ctxHandle))
}

//export swift_io_v_v23_context_CancelableVContext_nativeCancelAsync
func swift_io_v_v23_context_CancelableVContext_nativeCancelAsync(ctxHandle C.GoCancelableContextHandle, asyncId C.AsyncCallbackIdentifier, successCallback C.SwiftAsyncSuccessCallback) {
	ctx, cancelFunc := GoCancelableContext(uint64(ctxHandle))
	go func() {
		cancelFunc()
		<-ctx.Done()
		sutil.DoSuccessCallback(unsafe.Pointer(successCallback), int32(asyncId))
	}()
}
