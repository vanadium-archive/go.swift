package util

import (
	"math"
	"time"
)

// #import "types.h"
import "C"

// Useful when needing to return SOMETHING for a given function that otherwise is throwing an error via the errPtr
func EmptySwiftByteArray() C.SwiftByteArray {
	var empty C.SwiftByteArray
	empty.length = 0
	empty.data = nil
	return empty
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
