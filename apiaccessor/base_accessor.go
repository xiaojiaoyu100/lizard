package apiaccessor

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type baseAccessor struct {
	evalSignatureFunc EvalSignature
	timestampChecker  TimestampChecker
	nonceChecker      NonceChecker
	args              args
}

func defEvalSignatureFunc(move uint) EvalSignature {
	return func(origin string) (signature string) {
		md5Hash := md5.New()
		_, _ = md5Hash.Write([]byte(origin))
		checksumText := strings.ToLower(hex.EncodeToString(md5Hash.Sum(nil)))
		return checksumText[move:] + checksumText[:move]
	}
}

func defTimestampChecker(timestamp int64) error {
	const (
		sec        = 5
		timeFormat = "2006/01/02 15:04:05"
	)
	now := time.Now()
	dt := now.Unix() - timestamp
	if dt > sec || dt < -sec {
		nowTimeStr := now.Format(timeFormat)
		timestampStr := time.Unix(timestamp, 0).In(now.Location()).Format(timeFormat)
		return fmt.Errorf("%w: now %s, get %s", ErrTimestampTimeout, nowTimeStr, timestampStr)
	}
	return nil
}

func defNonceChecker(_ string) error {
	return nil
}

func newBaseAccessor() baseAccessor {
	const baseSignatureMoveSep = 2 // 签名md5移动2位
	return baseAccessor{
		evalSignatureFunc: defEvalSignatureFunc(baseSignatureMoveSep),
		timestampChecker:  defTimestampChecker,
		nonceChecker:      defNonceChecker,
		args:              newArgs(),
	}
}

// CheckSignature implements the Accessor CheckSignature interface
func (a *baseAccessor) CheckSignature() error {
	// 参数排序
	sort.Slice(a.args.l, func(i, j int) bool {
		return a.args.l[i].k < a.args.l[j].k
	})
	// 拼接参数key-value
	var (
		argText string
		i       int
	)
	for _, arg := range a.args.l {
		if arg.k == signatureTag {
			continue
		}
		if i == 0 {
			argText = fmt.Sprintf("%s=%s", arg.k, arg.v)
		} else {
			argText = fmt.Sprintf("%s&%s=%s", argText, arg.k, arg.v)
		}
		i++
	}

	// 比较签名
	signature := a.evalSignatureFunc(argText)
	argSignature := a.args.kv[signatureTag]
	if signature != argSignature {
		return fmt.Errorf("%w: want %s, get %s", ErrSignatureUnmatched, signature, argSignature)
	}
	return nil
}

// CheckTimestamp implements the Accessor CheckTimestamp interface
func (a *baseAccessor) CheckTimestamp() error {
	timestampStr := a.args.kv[timestampTag]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}
	return a.timestampChecker(timestamp)
}

// CheckNonce implements the Accessor CheckNonce interface
func (a *baseAccessor) CheckNonce() error {
	return a.nonceChecker(a.args.kv[nonceTag])
}
