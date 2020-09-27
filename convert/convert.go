package convert

import (
	"unsafe"
)

// stringHeader is the runtime representation of a string.
// It should be identical to reflect.StringHeader
type stringHeader struct {
	data      unsafe.Pointer
	stringLen int
}

// sliceHeader is the runtime representation of a slice.
// It should be identical to reflect.sliceHeader
type sliceHeader struct {
	data     unsafe.Pointer
	sliceLen int
	sliceCap int
}

// StringToByte unsafely converts s into a byte slice.
// If you modify b, then s will also be modified. This violates the
// property that strings are immutable.
func StringToByte(s string) (b []byte) {
	stringHeader := (*stringHeader)(unsafe.Pointer(&s))
	sliceHeader := (*sliceHeader)(unsafe.Pointer(&b))
	sliceHeader.data = stringHeader.data
	sliceHeader.sliceLen = len(s)
	sliceHeader.sliceCap = len(s)
	return b
}

// ByteToString unsafely converts b into a string.
// If you modify b, then s will also be modified. This violates the
// property that strings are immutable.
func ByteToString(b []byte) (s string) {
	sliceHeader := (*sliceHeader)(unsafe.Pointer(&b))
	stringHeader := (*stringHeader)(unsafe.Pointer(&s))
	stringHeader.data = sliceHeader.data
	stringHeader.stringLen = len(b)
	return s
}
