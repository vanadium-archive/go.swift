// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package v23

import (
	_ "v.io/x/swift/v23/context"
	_ "v.io/x/swift/v23/security"
)

import "C"

func Init() error {
	// Currently nothing needed. Placeholder for the potential future.
	return nil
}
