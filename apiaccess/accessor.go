package apiaccess

import (
	"errors"
)

type Accessor interface {
	CheckSignature() error
	CheckTimestamp() error
	CheckNonce() error
}

var (
	ErrArgLack            = errors.New("arg lack")
	ErrSignatureUnmatched = errors.New("signature is unmatched")
	ErrTimestampTimeout   = errors.New("timestamp time out")
	ErrNonceUsed          = errors.New("nonce is used")
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

type TimestampChecker func(timestamp int64) error

type NonceChecker func(nonce string) error
