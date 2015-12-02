// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin ios

package util

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

import "C"

// GoRef creates a new reference to the value addressed by the provided pointer.
// The value will remain referenced until it is explicitly unreferenced using
// goUnref() using the returned sequence id.
func GoNewRef(valptr interface{}) uint64 {
	if !IsPointer(valptr) {
		panic(fmt.Sprintf("must pass pointer value to goRef; instead got %v", valptr))
	}
	return goRefs.newRef(valptr)
}

// Increments the reference count to a previously existing id pointing to a value.
func GoRef(id uint64) {
	goRefs.ref(id)
}

// Removes a previously added reference to the value addressed by the
// sequence id.  If the value hasn't been ref-ed (a bug?), this unref will
// be a no-op.
func GoUnref(id uint64) {
	goRefs.unref(id)
}

func GoGetRef(id uint64) interface{} {
	return goRefs.get(id)
}

// Returns true iff the provided value is a pointer.
func IsPointer(val interface{}) bool {
	v := reflect.ValueOf(val)
	return v.Kind() == reflect.Ptr || v.Kind() == reflect.UnsafePointer
}

// Return the value of the pointer as a uintptr.
func PtrValue(ptr interface{}) uintptr {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.UnsafePointer {
		panic(fmt.Sprintf("must pass pointer value to PtrValue, was %v ", v.Type()))
	}
	return v.Pointer()
}

// Dereferences the provided (pointer) value, or panic-s if the value isn't of pointer type.
func DerefOrDie(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("want reflect.Ptr value for %v, have %v", i, v.Type()))
	}
	return v.Elem().Interface()
}

// goRefs stores references to instances of various Go types, namely instances
// that are referenced only by the Swift code.  The only purpose of this store
// is to prevent Go runtime from garbage collecting those instances.
var goRefs = newSafeRefCounter()
var lastId uint64 = 0

type refData struct {
	instance interface{}
	count    int
}

// Returns a new instance of a thread-safe reference counter.
func newSafeRefCounter() *safeRefCounter {
	return &safeRefCounter{
		refs: make(map[uint64]*refData),
	}
}

// safeRefCounter is a thread-safe reference counter.
type safeRefCounter struct {
	lock sync.Mutex
	refs map[uint64]*refData
}

// Increment the reference count to the given valptr.
func (c *safeRefCounter) ref(id uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	ref, ok := c.refs[id]
	if !ok {
		panic(fmt.Sprintf("Refing id %d that doesn't exist", id))
	} else {
		ref.count++
	}
}

// Given a valptr, store that into our map and return the associated handle for Swift
func (c *safeRefCounter) newRef(valptr interface{}) uint64 {
	//	p := PtrValue(valptr)
	c.lock.Lock()
	defer c.lock.Unlock()
	id := atomic.AddUint64(&lastId, 1)
	c.refs[id] = &refData{
		instance: valptr,
		count:    1,
	}
	return id
}

// Decrement the reference count of the valptr asssociated with the handle, returning
// the new reference count value and deleting the assocation if it hits 0.
func (c *safeRefCounter) unref(id uint64) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	ref, ok := c.refs[id]
	if !ok {
		panic(fmt.Sprintf("Unrefing id %d that hasn't been refed before", id))
	}
	count := ref.count
	if count == 0 {
		panic(fmt.Sprintf("Ref count for id %d is zero", id))
	}
	if count > 1 {
		ref.count--
		return ref.count
	}
	delete(c.refs, id)
	return 0
}

// Given a handle, return the associated go object
func (c *safeRefCounter) get(id uint64) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	ref, ok := c.refs[id]
	if !ok {
		panic(fmt.Sprintf("Trying to get id %d that doesn't exist", id))
	} else {
		return ref.instance
	}
}
