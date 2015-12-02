// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package util

import (
	"unsafe"

	"v.io/v23/context"
	"v.io/v23/verror"
)

/*
#include "types.h"

// Because of weirdness with function pointers & cgo, we hack around limitations by defining methods that just call
// C function pointers here. Go can call these functions directly.

typedef void (*fxn_AsyncId)(AsyncCallbackIdentifier);
typedef void (*fxn_AsyncId_GoHandle)(AsyncCallbackIdentifier, _GoHandle);
typedef void (*fxn_AsyncId_SwiftByteArrayArray)(AsyncCallbackIdentifier, SwiftByteArrayArray);
typedef void (*fxn_AsyncId_SwiftVError)(AsyncCallbackIdentifier, SwiftVError);

static void callFxn_AsyncId(fxn_AsyncId fxnPtr, AsyncCallbackIdentifier a) {
	fxnPtr(a);
}

static void callFxn_AsyncId_GoHandle(fxn_AsyncId_GoHandle fxnPtr, AsyncCallbackIdentifier a, _GoHandle b) {
	fxnPtr(a, b);
}

static void callFxn_AsyncId_SwiftByteArrayArray(
		fxn_AsyncId_SwiftByteArrayArray fxnPtr, AsyncCallbackIdentifier a, SwiftByteArrayArray b) {
	fxnPtr(a, b);
}

static void callFxn_AsyncId_SwiftVError(fxn_AsyncId_SwiftVError fxnPtr, AsyncCallbackIdentifier a, SwiftVError b) {
	fxnPtr(a, b);
}
*/
import "C"

func ThrowSwiftError(ctx *context.T, err error, swiftVErrorStructPtr unsafe.Pointer) {
	id := verror.ErrorID(err)
	actionCode := verror.Action(err)
	vErr := verror.Convert(verror.IDAction{id, actionCode}, ctx, err)
	pcs := verror.Stack(vErr)
	stacktrace := pcs.String()
	msg := vErr.Error()

	var swiftErrorPtr *C.SwiftVError = (*C.SwiftVError)(swiftVErrorStructPtr)
	(*swiftErrorPtr).identity = C.CString((string)(id))
	(*swiftErrorPtr).actionCode = C._GoUint32(actionCode)
	(*swiftErrorPtr).msg = C.CString(msg)
	(*swiftErrorPtr).stacktrace = C.CString(stacktrace)
}

func DoSuccessCallback(successPtr unsafe.Pointer, asyncId int32) {
	C.callFxn_AsyncId(
		(C.fxn_AsyncId)(successPtr),
		(C.AsyncCallbackIdentifier)(asyncId))
}

func DoSuccessHandlerCallback(successPtr unsafe.Pointer, asyncId int32, handler uint64) {
	C.callFxn_AsyncId_GoHandle(
		(C.fxn_AsyncId_GoHandle)(successPtr),
		(C.AsyncCallbackIdentifier)(asyncId),
		(C._GoHandle)(handler))
}

func DoSuccessByteArrayArrayCallback(successPtr unsafe.Pointer, asyncId int32, swiftByteArrayArrayPtr unsafe.Pointer) {
	var byteArrayArray C.SwiftByteArrayArray = *(*C.SwiftByteArrayArray)(swiftByteArrayArrayPtr)
	C.callFxn_AsyncId_SwiftByteArrayArray(
		(C.fxn_AsyncId_SwiftByteArrayArray)(successPtr),
		(C.AsyncCallbackIdentifier)(asyncId),
		(C.SwiftByteArrayArray)(byteArrayArray))
}

func DoFailureCallback(failurePtr unsafe.Pointer, asyncId int32, swiftVErrorStructPtr unsafe.Pointer) {
	var swiftError C.SwiftVError = *(*C.SwiftVError)(swiftVErrorStructPtr)
	C.callFxn_AsyncId_SwiftVError(
		(C.fxn_AsyncId_SwiftVError)(failurePtr),
		(C.AsyncCallbackIdentifier)(asyncId),
		swiftError)
}
