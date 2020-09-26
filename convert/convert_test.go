package convert_test

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/xiaojiaoyu100/lizard/convert"
)

func TestString2Byte(t *testing.T) {
	s := "1234567890"
	b := convert.String2Byte(s)
	assert.Equal(t, len(b), len(s))
	assert.Equal(t, cap(b), len(s))
}

func TestByteToString(t *testing.T) {
	b := []byte("1234567890")
	s := convert.ByteToString(b)
	assert.Equal(t, len(b), len(s))
	assert.Equal(t, cap(b), len(s))
}
