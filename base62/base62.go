package base62

import (
	"strings"

	"github.com/xiaojiaoyu100/lizard/stringkit"
)

const (
	encodeScheme = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Encode convert a int64 number to Base62 string.
func Encode(n int64) string {
	if n == 0 {
		return "0"
	}
	var ret string
	for n != 0 {
		r := n % 62
		ret += string(encodeScheme[r])
		n /= 62
	}
	return stringkit.Reverse(ret)
}

// Decode converts a Base62 string to a int64 number.
func Decode(s string) int64 {
	l := len(s)
	var ret int64
	for i := 0; i < l; i++ {
		ret = 62*ret + int64(strings.Index(encodeScheme, string(s[i])))
	}
	return ret
}
