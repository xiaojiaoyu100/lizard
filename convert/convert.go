package convert

import "unsafe"

// ByteToString converts a byte slice to string.
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Byte converts a string to byte slice.
func String2Byte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
