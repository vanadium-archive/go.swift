// Copyright 2016 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package discovery

import (
	"encoding/json"
	"unsafe"

	"v.io/v23"
	"v.io/v23/discovery"
	"v.io/v23/security"
	sutil "v.io/x/swift/util"
	scontext "v.io/x/swift/v23/context"
)

/*
#include <stdlib.h>
#import "../../types.h"

static void CallAdvertisingCallback(SwiftAsyncSuccessCallback callback, AsyncCallbackIdentifier asyncId) {
  callback(asyncId);
}
static void CallScanCallback(SwiftAsyncJsonCallback callback, AsyncCallbackIdentifier asyncId, SwiftByteArray data) {
  callback(asyncId, data);
}
*/
import "C"

//export swift_io_v_v23_discovery_new
func swift_io_v_v23_discovery_new(ctxHandle C.GoContextHandle, errOut *C.SwiftVError) C.GoDiscoveryHandle {
	ctx := scontext.GoContext(uint64(ctxHandle))
	d, err := v23.NewDiscovery(ctx)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return C.GoDiscoveryHandle(0)
	}
	return C.GoDiscoveryHandle(sutil.GoNewRef(&d))
}

//export swift_io_v_v23_discovery_finalize
func swift_io_v_v23_discovery_finalize(discoveryHandle C.GoDiscoveryHandle) {
	sutil.GoUnref(uint64(discoveryHandle))
}

// Exports the discovery advertise API to CGO using JSON to marshal ads
//export swift_io_v_v23_discovery_advertise
func swift_io_v_v23_discovery_advertise(ctxHandle C.GoContextHandle, discoveryHandle C.GoDiscoveryHandle, adJson C.SwiftByteArray, visibilityArray C.SwiftCStringArray, asyncId C.AsyncCallbackIdentifier, doneCallback C.SwiftAsyncSuccessCallback, errOut *C.SwiftVError) bool {
	ctx := scontext.GoContext(uint64(ctxHandle))
	d := GoDiscoveryT(uint64(discoveryHandle))
	ad := discovery.Advertisement{}
	if err := json.Unmarshal(sutil.GoBytesNoCopy(unsafe.Pointer(&adJson)), &ad); err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return false
	}
	var visibility []security.BlessingPattern
	for _, v := range visibilityArray.toStrings() {
		visibility = append(visibility, security.BlessingPattern(v))
	}
	doneChan, err := d.Advertise(ctx, &ad, visibility)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return false
	}
	go func() {
		<-doneChan
		C.CallAdvertisingCallback(doneCallback, asyncId)
	}()
	return true
}

// Exports the discovery scan API to CGO using JSON to marshal ads
//export swift_io_v_v23_discovery_scan
func swift_io_v_v23_discovery_scan(ctxHandle C.GoContextHandle, discoveryHandle C.GoDiscoveryHandle, query C.SwiftCString, asyncId C.AsyncCallbackIdentifier, callbackBlock C.SwiftAsyncJsonCallback, errOut *C.SwiftVError) bool {
	ctx := scontext.GoContext(uint64(ctxHandle))
	d := GoDiscoveryT(uint64(discoveryHandle))
	goQuery := sutil.GoString(unsafe.Pointer(&query), false)
	ch, err := d.Scan(ctx, goQuery)
	if err != nil {
		sutil.ThrowSwiftError(ctx, err, unsafe.Pointer(errOut))
		return false
	}
	go func() {
		for update := range ch {
			data := struct {
				IsLost bool
				Ad     discovery.Advertisement
			}{
				IsLost: update.IsLost(),
				Ad:     update.Advertisement(),
			}
			b, err := json.Marshal(data)
			if err != nil {
				ctx.Fatal("Unable to JSON serialize discovery update: ", err)
			}
			ba := swiftBytesCopy(b)
			C.CallScanCallback(callbackBlock, asyncId, ba)
			C.free(ba.data)
		}
		C.CallScanCallback(callbackBlock, asyncId, emptySwiftByteArray())
	}()
	return true
}
