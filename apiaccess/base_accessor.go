package apiaccess

import (
	"fmt"
	"sort"
	"strconv"
)

type baseAccessor struct {
	evalSignatureFunc EvalSignature
	evalNonceKeyFunc  EvalNonceKey
	timestampChecker  TimestampChecker
	nonceChecker      NonceChecker
	args              args
}

func newBaseAccessor() baseAccessor {
	const baseSignatureMoveSep = 2 // 签名md5移动2位
	return baseAccessor{
		evalSignatureFunc: defEvalSignatureFunc(baseSignatureMoveSep),
		evalNonceKeyFunc:  defEvalNonceKey,
		timestampChecker:  defTimestampChecker,
		nonceChecker:      defNonceChecker,
		args:              newArgs(),
	}
}

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
		return fmt.Errorf("%w: want %s, get %s", ErrSignatureUnmatch, signature, argSignature)
	}
	return nil
}

func (a *baseAccessor) CheckTimestamp() error {
	timestampStr := a.args.kv[timestampTag]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}
	return a.timestampChecker(timestamp)
}

func (a *baseAccessor) CheckNonce() error {
	return a.nonceChecker(a.args.kv[nonceTag])
}
