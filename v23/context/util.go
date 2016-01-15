// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package context

import (
	"fmt"

	"v.io/v23/context"

	sutil "v.io/x/swift/util"
)

// #include "../../types.h"
import "C"

// Converts a Go Context into a native pointer for Swift and increments the go reference count
func SwiftContext(ctx *context.T) C.GoContextHandle {
	id := sutil.GoNewRef(ctx)
	return C.GoContextHandle(id)
}

// Converts a native pointer from Swift to a Go Context. Panics if it can't.
func GoContext(ctxHandle uint64) *context.T {
	valptr := sutil.GoGetRef(ctxHandle)
	if ctx, ok := valptr.(*context.T); ok {
		return ctx
	} else {
		panic(fmt.Sprintf("Couldn't get context from handle with id %d", ctxHandle))
	}
}

type cancelFuncKey struct{}

func SwiftCancelableContext(ctx *context.T, cancelFunc context.CancelFunc) (C.GoCancelableContextHandle, error) {
	if cancelFunc == nil {
		return C.GoCancelableContextHandle(0), fmt.Errorf("Cannot create SwiftCancelableContext with nil cancel function")
	}
	ctx = context.WithValue(ctx, cancelFuncKey{}, cancelFunc)
	return C.GoCancelableContextHandle(SwiftContext(ctx)), nil
}

// Converts a native pointer from Swift to a Go Cancelable Context. Panics if it can't.
func GoCancelableContext(ctxHandle uint64) (*context.T, context.CancelFunc) {
	ctx := GoContext(ctxHandle)
	value := ctx.Value(cancelFuncKey{})
	if cancelFunc, ok := value.(context.CancelFunc); ok {
		return ctx, cancelFunc
	} else {
		panic(fmt.Sprintf("Couldn't cast cancelFunc for contextId %d", ctxHandle))
	}
}
