package convert

import "unsafe"

// ByteToString converts a byte slice to string.
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToByte converts a string to byte slice.
func StringToByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
