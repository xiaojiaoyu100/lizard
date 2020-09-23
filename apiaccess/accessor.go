package apiaccess

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/xiaojiaoyu100/lizard/convert"
	"strings"
	"time"
)

type Accessor interface {
	// CheckSignature compare the argument signature with another which send from client.
	// if they are un match, CheckSignature will return the error, errSignatureUnmatch.
	CheckSignature() error
	CheckTimestamp() error
	CheckNonce() error
}

var (
	ErrArgLack          = errors.New("arg lack")
	ErrSignatureUnmatch = errors.New("signature unmatch")
	ErrTimestampTimeout = errors.New("timestamp time out")
	ErrNonceUsed        = errors.New("nonce is used")
)

const (
	nonceTag     = "nonce"
	timestampTag = "timestamp"
	secretKeyTag = "secret_key"
	signatureTag = "signature"
)

type arg struct {
	k string
	v string
}

type args struct {
	kv map[string]string
	l  []*arg
}

func newArgs() args {
	return args{
		kv: make(map[string]string),
		l:  make([]*arg, 0),
	}
}

func (a *args) append(k, v string) {
	a.kv[k] = v
	a.l = append(a.l, &arg{k: k, v: v})
}

type EvalSignature func(origin string) (signature string)

type EvalNonceKey func(nonce string) (key string)

type TimestampChecker func(timestamp int64) error

type NonceChecker func(nonce string) error

func defEvalSignatureFunc(move uint) EvalSignature {
	return func(origin string) (signature string) {
		md5Hash := md5.New()
		_, _ = md5Hash.Write(convert.String2Byte(origin))
		checksumText := strings.ToLower(hex.EncodeToString(md5Hash.Sum(nil)))
		return checksumText[move:] + checksumText[:move]
	}
}

func defEvalNonceKey(nonce string) (key string) {
	return nonce
}

func defTimestampChecker(timestamp int64) error {
	const sec = 5
	now := time.Now().Unix()
	if now-timestamp > sec || timestamp-now > sec {
		return ErrTimestampTimeout
	}
	return nil
}

func defNonceChecker(nonce string) error {
	return nil
}
