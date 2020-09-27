package apiaccessor

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestNewQueryAccessor(t *testing.T) {
	query := url.Values{}
	_, err := NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, errArgLack), true)

	query = url.Values{
		nonceTag: []string{"12345"},
	}
	_, err = NewQueryAccessor(query, "123")
	assert.Equal(t, errors.Is(err, errArgLack), true)

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
	assert.Equal(t, errors.Is(err, errSignatureUnmatched), true)

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

func TestCheckTimestamp(t *testing.T) {
	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckTimestamp()
	assert.Equal(t, errors.Is(err, errTimestampTimeout), true)

	query = url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err = NewQueryAccessor(query, "123")
	assert.Equal(t, err, nil)
	err = qa.CheckTimestamp()
	assert.Equal(t, err, nil)
}

func TestCheckNonce(t *testing.T) {
	nonceMap := make(map[string]bool)
	mockNonceChecker := func(nonce string) error {
		if _, ok := nonceMap[nonce]; ok {
			return errNonceUsed
		}
		nonceMap[nonce] = true
		return nil
	}

	query := url.Values{
		nonceTag:     []string{"12345"},
		signatureTag: []string{"ca444a9db0301178257b0d9e959533a3"},
		timestampTag: []string{"12345"},
		"phone":      []string{"12345"},
		"abc":        []string{"abc"},
	}
	qa, err := NewQueryAccessor(query, "123", WithNonceChecker(mockNonceChecker))
	assert.Equal(t, err, nil)
	err = qa.CheckNonce()
	assert.Equal(t, err, nil)
	err = qa.CheckNonce()
	assert.Equal(t, errors.Is(err, errNonceUsed), true)
}
