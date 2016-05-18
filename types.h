// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// +build darwin,ios
//
// This file defines the types and structs that get passed between the Go runtime and Swift

// Match the generated types from CGo without being able to include it's generated conversions
typedef unsigned long long _GoUint64;
typedef unsigned int _GoUint32;
typedef long long _GoInt64;
typedef int _GoInt32;

// Define the handle types that allow us to reference objects in Go from Swift without passing unsafe pointers
// This is particularly important given the upcoming moving GC dictates we can't rely on pointers staying stable.
// See: https://github.com/golang/proposal/blob/master/design/12416-cgo-pointers.md
typedef _GoUint64 _GoHandle;
typedef _GoHandle GoContextHandle;
typedef _GoHandle GoCancelableContextHandle;
typedef _GoHandle GoClientCallHandle;
typedef _GoHandle GoDiscoveryHandle;

typedef struct {
    _GoUint64 length;
    char *data;
} SwiftCString;

typedef struct {
    _GoUint64 length;
    SwiftCString *data;
} SwiftCStringArray;

// Express byte arrays with lengths
typedef struct {
    _GoUint64 length;
    void* data;
} SwiftByteArray;

// Express arrays of byte arrays with lengths
typedef struct {
    _GoUint64 length;
    SwiftByteArray* data;
} SwiftByteArrayArray;

// Encodes Go's VError in a format suitable for Swift conversion
typedef struct {
	char* identity;            // The identity of the error.
	_GoUint32 actionCode;      // Default action to take on error.
	char* msg;                 // Error message; empty if no language known.
	char* stacktrace;          // Stacktraces rendered into a single string
} SwiftVError;
typedef SwiftVError* SwiftVErrorPtr;

// Asynchronous callback function pointers, which use their own handle (AsyncCallbackIdentifier) to overcome
// Swift limitations on context-free closures (function pointers)
typedef int AsyncCallbackIdentifier;
typedef void (*SwiftAsyncJsonCallback)(AsyncCallbackIdentifier asyncId, SwiftByteArray json);
typedef void (*SwiftAsyncSuccessCallback)(AsyncCallbackIdentifier asyncId);
typedef void (*SwiftAsyncSuccessHandleCallback)(AsyncCallbackIdentifier asyncId, _GoHandle handle);
typedef void (*SwiftAsyncSuccessByteArrayArrayCallback)(AsyncCallbackIdentifier asyncId, SwiftByteArrayArray byteArrayArray);
typedef void (*SwiftAsyncFailureCallback)(AsyncCallbackIdentifier asyncId, SwiftVError err);

// Other callbacks
typedef void (*SwiftAsyncVoidCallback)();

