// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package google

import (
	_ "v.io/x/ios/impl/google/rpc"
	_ "v.io/x/ios/impl/google/rt"
)

func Init() error {
	// Currently nothing needed. Placeholder for the potential future.
	return nil
}
