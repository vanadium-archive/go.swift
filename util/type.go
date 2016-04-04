// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package util

import (
	"encoding/base64"
	"math"
	"time"
	"unsafe"

	"v.io/x/lib/vlog"
)

/*
#include <stdlib.h>
#import "../types.h"
*/
import "C"

// Useful when needing to return SOMETHING for a given function that otherwise is throwing an error via the errPtr
func EmptySwiftByteArray() C.SwiftByteArray {
	var empty C.SwiftByteArray
	empty.length = 0
	empty.data = nil
	return empty
}

func GoBytesCopy(swiftByteArrayPtr unsafe.Pointer) []byte {
	swiftArray := *(*C.SwiftByteArray)(swiftByteArrayPtr)
	length := C.int(swiftArray.length)
	return C.GoBytes(swiftArray.data, length)
}

func GoBytesNoCopy(swiftByteArrayPtr unsafe.Pointer) []byte {
	// Taken from https://github.com/golang/go/wiki/cgo
	// "To create a Go slice backed by a C array (without copying the original data),
	// one needs to acquire this length at runtime and use a type conversion to a pointer
	// to a very big array and then slice it to the length that you want"
	swiftArray := *(*C.SwiftByteArray)(swiftByteArrayPtr)
	length := int(swiftArray.length)
	var bytes []byte = (*[1 << 30]byte)(unsafe.Pointer(swiftArray.data))[:length:length]
	return bytes
}

// Utils to convert between Go times and durations and NSTimeInterval in Swift
func NSTimeInterval(t time.Time) C.double {
	return C.double(t.UnixNano() / 1000000000.0)
}

func GoTime(t float64) time.Time {
	seconds := math.Floor(t)
	nsec := (t - seconds) * 1000000000.0
	return time.Unix(int64(seconds), int64(nsec))
}

func GoDuration(d float64) time.Duration {
	return time.Duration(int64(d * 1e9))
}

//export swift_io_v_swift_impl_util_type_nativeBase64UrlDecode
func swift_io_v_swift_impl_util_type_nativeBase64UrlDecode(base64UrlEncoded *C.char) C.SwiftByteArray {
	// Decode the base64 url encoded string to bytes in a way that prevents extra copies along the CGO boundary.
	urlEncoded := C.GoString(base64UrlEncoded)
	maxLength := base64.URLEncoding.DecodedLen(len(urlEncoded))
	bytesBacking := C.malloc(C.size_t(maxLength))
	if bytesBacking == nil {
		vlog.Errorf("Unable allocate %v bytes", maxLength)
		return EmptySwiftByteArray()
	}
	var bytes []byte = (*[1 << 30]byte)(unsafe.Pointer(bytesBacking))[:maxLength:maxLength]
	n, err := base64.URLEncoding.Decode(bytes, []byte(urlEncoded))
	if err != nil {
		vlog.Errorf("Unable to base64 decode string: %v\n", err)
		C.free(bytesBacking)
		return EmptySwiftByteArray()
	}
	var swiftArray C.SwiftByteArray
	swiftArray.length = C._GoUint64(n)
	swiftArray.data = bytesBacking
	return swiftArray
}
