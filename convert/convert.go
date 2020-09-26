package convert

import (
	"reflect"
	"unsafe"
)

// ByteToString converts a byte slice to string.
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Byte converts a string to byte slice.
func String2Byte(s string) (b []byte) {
	b = *(*[]byte)(unsafe.Pointer(&s))
	(*reflect.SliceHeader)(unsafe.Pointer(&b)).Cap = len(s)
	return
}
