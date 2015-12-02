// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package ios

import (
	"unsafe"

	"v.io/x/lib/vlog"

	// TODO Make this pluggable somehow
	_ "v.io/x/ref/runtime/factories/roaming"

	igoogle "v.io/x/ios/impl/google"
	iutil "v.io/x/ios/util"
	iv23 "v.io/x/ios/v23"
)

//#import "types.h"
import "C"

//export ios_io_v_v23_V_nativeInitGlobal
func ios_io_v_v23_V_nativeInitGlobal(errOut *C.SwiftVError) {
	// Send all vlog logs to stderr during the init so that we don't crash on android trying
	// to create a log file.  These settings will be overwritten in nativeInitLogging below.
	vlog.Log.Configure(vlog.OverridePriorConfiguration(true), vlog.LogToStderr(true))

	if err := iv23.Init(); err != nil {
		iutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
	if err := igoogle.Init(); err != nil {
		iutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
}

//export ios_io_v_v23_V_nativeInitLogging
func ios_io_v_v23_V_nativeInitLogging(logDir *C.char, logToStderr bool, logLevel int, moduleSpec *C.char, errOut *C.SwiftVError) {
	dir, toStderr, level, vmodule, err := loggingOpts(C.GoString(logDir), logToStderr, logLevel, C.GoString(moduleSpec))
	if err != nil {
		iutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}

	vlog.Log.Configure(vlog.OverridePriorConfiguration(true), dir, toStderr, level, vmodule)
}

func loggingOpts(logDir string, logToStderr bool, logLevel int, moduleSpec string) (dir vlog.LogDir, toStderr vlog.LogToStderr, level vlog.Level, vmodule vlog.ModuleSpec, err error) {
	dir = vlog.LogDir(logDir)
	toStderr = vlog.LogToStderr(logToStderr)
	level = vlog.Level(logLevel)
	err = vmodule.Set(moduleSpec)
	return
}
