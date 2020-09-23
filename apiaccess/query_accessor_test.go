package apiaccess

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestNewQueryAccessor(t *testing.T) {
	query := url.Values{}
	_, err := NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, ErrArgLack), true)

	query = url.Values{
		nonceTag: []string{"12345"},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, ErrArgLack), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"a":          []string{"12345"},
		"b":          []string{"12345"},
		"c":          []string{"12345"},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
}

func TestCheckSignature(t *testing.T) {
	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"12345"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckSignature()
	assert.Equal(t, errors.Is(err, ErrSignatureUnmatch), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckSignature()
	assert.Equal(t, err, nil)
}
