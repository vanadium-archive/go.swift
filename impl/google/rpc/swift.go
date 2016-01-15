// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package rpc

import (
	"reflect"
	"unsafe"

	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/options"
	"v.io/v23/rpc"
	"v.io/v23/security"
	"v.io/v23/vdl"
	"v.io/v23/vom"

	sutil "v.io/x/swift/util"
	scontext "v.io/x/swift/v23/context"
)

/*
#include <string.h> // memcpy
#import "../../../types.h"
// These sizes (including C struct memory alignment/padding) isn't available from Go, so we make that available via CGo.
static const size_t sizeofSwiftByteArray = sizeof(SwiftByteArray);
static const size_t sizeofSwiftByteArrayArray = sizeof(SwiftByteArrayArray);
*/
import "C"

func doStartCall(context *context.T, name, method string, skipServerAuth bool, client rpc.Client, args []interface{}) (rpc.ClientCall, error) {
	var opts []rpc.CallOpt
	if skipServerAuth {
		opts = append(opts,
			options.NameResolutionAuthorizer{security.AllowEveryone()},
			options.ServerAuthorizer{security.AllowEveryone()})
	}
	// Invoke StartCall
	call, err := client.StartCall(context, name, method, args, opts...)
	if err != nil {
		return nil, err
	}
	return call, nil
}

// TODO From JNI, implement this for Swift
//func decodeArgs(cVomArgs C.SwiftByteArrayArray) ([]interface{}, error) {
//	if (cVomArgs == nil || cVomArgs.length == 0) {
//		return make([]interface{}, 0), nil
//	}
//	// VOM-decode each arguments into a *vdl.Value.
//	args := make([]interface{}, cVomArgs.length)
//	for i := 0; i < cVomArgs.length; i++ {
//		ret[i] = byte(*ptr)
//		ptr = (*C.jbyte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*ptr)))
//
//		cVomArgs
//
//		var err error
//		if args[i], err = jutil.VomDecodeToValue(vomArgs[i]); err != nil {
//			return nil, err
//		}
//	}
//	return args, nil
//}

//export swift_io_v_impl_google_rpc_ClientImpl_nativeStartCallAsync
func swift_io_v_impl_google_rpc_ClientImpl_nativeStartCallAsync(ctxHandle C.GoContextHandle, cName *C.char, cMethod *C.char, cVomArgs C.SwiftByteArrayArray, skipServerAuth bool, asyncId C.AsyncCallbackIdentifier, successCallback C.SwiftAsyncSuccessHandleCallback, failureCallback C.SwiftAsyncFailureCallback) {
	name := C.GoString(cName)
	method := C.GoString(cMethod)
	ctx := scontext.GoContext(uint64(ctxHandle))
	client := v23.GetClient(ctx)

	// TODO Get args (we don't have VOM yet in Swift so nothing to get until then)
	//	args, err := decodeArgs(env, jVomArgs)
	//	if err != nil {
	//		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
	//		return C.GoClientCallHandle(0)
	//	}
	args := make([]interface{}, 0)

	go func() {
		result, err := doStartCall(ctx, name, method, skipServerAuth == true, client, args)
		if err != nil {
			var swiftVError C.SwiftVError
			sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(&swiftVError))
			sutil.DoFailureCallback(unsafe.Pointer(failureCallback), int32(asyncId), unsafe.Pointer(&swiftVError))
		} else {
			handle := C.GoClientCallHandle(SwiftClientCall(result))
			sutil.DoSuccessHandlerCallback(unsafe.Pointer(successCallback), int32(asyncId), uint64(handle))
		}
	}()
}

//export swift_io_v_impl_google_rpc_ClientImpl_nativeClose
func swift_io_v_impl_google_rpc_ClientImpl_nativeClose(ctxHandle C.GoContextHandle) {
	ctx := scontext.GoContext(uint64(ctxHandle))
	client := v23.GetClient(ctx)
	<-client.Closed()
}

//export swift_io_v_impl_google_rpc_ClientCallImpl_nativeCloseSend
func swift_io_v_impl_google_rpc_ClientCallImpl_nativeCloseSend(ctxHandle C.GoContextHandle, callHandle C.GoClientCallHandle, errOut *C.SwiftVError) {
	ctx := scontext.GoContext(uint64(ctxHandle))
	call := GoClientCall(uint64(callHandle))
	if err := call.CloseSend(); err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return
	}
}

func doFinish(call rpc.ClientCall, numResults int) (C.SwiftByteArrayArray, error) {
	// Have all the results be decoded into *vdl.Value.
	resultPtrs := make([]interface{}, numResults)
	for i := 0; i < numResults; i++ {
		value := new(vdl.Value)
		resultPtrs[i] = &value
	}
	if err := call.Finish(resultPtrs...); err != nil {
		// Invocation error.
		return EmptySwiftByteArrayArray(), err
	}

	// VOM-encode the results. Note in the future we'll want a pathway where we can get the original VOM results
	// from finish so we don't end up wasting CPU & memory here.

	// Prepare the byte array array that can be accessed from Swift via C.malloc
	vomResultsMemory := C.malloc(C.size_t(numResults * int(C.sizeofSwiftByteArray)))
	// Make that malloc'd memory available as a slice to Go.
	vomResultsPtrsHdr := reflect.SliceHeader{
		Data: uintptr(vomResultsMemory),
		Len:  numResults,
		Cap:  numResults,
	}
	vomResults := *(*[]C.SwiftByteArray)(unsafe.Pointer(&vomResultsPtrsHdr))
	// Create the C Struct to return that encapsulates our byte array array
	var cVomResults C.SwiftByteArrayArray
	cVomResults.length = C._GoUint64(numResults)
	cVomResults.data = (*C.SwiftByteArray)(vomResultsMemory)

	// For each result, VOM encode into a byte array that we stick into the returned struct
	for i, resultPtr := range resultPtrs {
		// Remove the pointer from the result.  Simply *resultPtr doesn't work
		// as resultPtr is of type interface{}.
		result := interface{}(sutil.DerefOrDie(resultPtr))
		var vomResult []byte
		var err error
		if vomResult, err = vom.Encode(result); err != nil {
			return EmptySwiftByteArrayArray(), err
		}
		cVomResultCopy := C.malloc(C.size_t(len(vomResult)))
		C.memcpy(cVomResultCopy, unsafe.Pointer(&vomResult[0]), C.size_t(len(vomResult)))
		var cVomResult C.SwiftByteArray
		cVomResult.length = C._GoUint64(len(vomResult))
		cVomResult.data = (*C.char)(cVomResultCopy)
		vomResults[i] = cVomResult
	}
	return cVomResults, nil
}

//export swift_io_v_impl_google_rpc_ClientCallImpl_nativeFinishAsync
func swift_io_v_impl_google_rpc_ClientCallImpl_nativeFinishAsync(ctxHandle C.GoContextHandle, callHandle C.GoClientCallHandle, numResults int, asyncId C.AsyncCallbackIdentifier, successCallback C.SwiftAsyncSuccessByteArrayArrayCallback, failureCallback C.SwiftAsyncFailureCallback) {
	ctx := scontext.GoContext(uint64(ctxHandle))
	call := GoClientCall(uint64(callHandle))
	go func() {
		result, err := doFinish(call, numResults)
		if err != nil {
			var swiftVError C.SwiftVError
			sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(&swiftVError))
			sutil.DoFailureCallback(unsafe.Pointer(failureCallback), int32(asyncId), unsafe.Pointer(&swiftVError))
		} else {
			sutil.DoSuccessByteArrayArrayCallback(unsafe.Pointer(successCallback), int32(asyncId), unsafe.Pointer(&result))
		}
	}()
}

//export swift_io_v_impl_google_rpc_ClientCallImpl_nativeFinalize
func swift_io_v_impl_google_rpc_ClientCallImpl_nativeFinalize(callHandle C.GoClientCallHandle) {
	sutil.GoUnref(uint64(callHandle))
}
