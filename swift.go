// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin,ios

package swift

import (
	"os"
	"unsafe"

	"v.io/x/lib/vlog"
	"v.io/x/ref"
	// TODO Make this pluggable somehow
	_ "v.io/x/ref/runtime/factories/roaming"
	sgoogle "v.io/x/swift/impl/google"
	sutil "v.io/x/swift/util"
	sv23 "v.io/x/swift/v23"
)

//#import "types.h"
import "C"

//export swift_io_v_v23_V_nativeInitGlobal
func swift_io_v_v23_V_nativeInitGlobal(credentialsDir *C.char, errOut *C.SwiftVError) {
	// Send all vlog logs to stderr during the init so that we don't crash on android trying
	// to create a log file.  These settings will be overwritten in nativeInitLogging below.
	vlog.Log.Configure(vlog.OverridePriorConfiguration(true), vlog.LogToStderr(true))

	if credentialsDir != nil {
		dir := C.GoString(credentialsDir)
		os.Setenv(ref.EnvCredentials, dir)
	}
	if err := sv23.Init(); err != nil {
		sutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
	if err := sgoogle.Init(); err != nil {
		sutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
		return
	}
}

//export swift_io_v_v23_V_nativeInitLogging
func swift_io_v_v23_V_nativeInitLogging(logDir *C.char, logToStderr bool, logLevel int, moduleSpec *C.char, errOut *C.SwiftVError) {
	dir, toStderr, level, vmodule, err := loggingOpts(C.GoString(logDir), logToStderr, logLevel, C.GoString(moduleSpec))
	if err != nil {
		sutil.ThrowSwiftError(nil, err, unsafe.Pointer(errOut))
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
