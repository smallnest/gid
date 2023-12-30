// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package goroutine provides functions that will return the runtime's
// ID number for the calling goroutine or its creator.
//
// The implementation is derived from Laevus Dexter's comment in Gophers' Slack #darkarts,
// https://gophers.slack.com/archives/C1C1YSQBT/p1593885226448300 post which linked to
// this playground snippet https://play.golang.org/p/CSOp9wyzydP.
package goroutine

import (
	"reflect"
	"unsafe"
)

// ID returns the runtime ID of the calling goroutine.
func ID() int64 {
	return idOf(getg(), goidoff)
}

func idOf(g unsafe.Pointer, off uintptr) int64 {
	return *(*int64)(add(g, off))
}

//go:nosplit
func getg() unsafe.Pointer {
	return *(*unsafe.Pointer)(add(getm(), curgoff))
}

type p struct {
	id int32
}

type puintptr uintptr

//go:nosplit
func (pp puintptr) ptr() *p { return (*p)(unsafe.Pointer(pp)) }

//go:nosplit
func (pp *puintptr) set(p *p) { *pp = puintptr(unsafe.Pointer(p)) }

// PID returns the "P" ID of the calling goroutine.
func PID() int32 {
	return getp().id
}

func pidOf(g unsafe.Pointer, off uintptr) int32 {
	return *(*int32)(add(g, off))
}

//go:nosplit
func getp() p {
	pp := (puintptr)(*(*unsafe.Pointer)(add(getm(), poff)))

	return *pp.ptr()
}

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname getm runtime.getm
//go:nosplit
func getm() unsafe.Pointer

var (
	curgoff = offset("*runtime.m", "curg")
	goidoff = offset("*runtime.g", "goid")

	poff = offset("*runtime.m", "p")
)

// offset returns the offset into typ for the given field.
func offset(typ, field string) uintptr {
	rt := toType(typesByString(typ)[0])
	f, _ := rt.Elem().FieldByName(field)
	return f.Offset
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

//go:linkname toType reflect.toType
func toType(t unsafe.Pointer) reflect.Type
