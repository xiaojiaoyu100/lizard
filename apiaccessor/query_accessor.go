package apiaccessor

import (
	"fmt"
	"net/url"
)

type QueryAccessor struct {
	baseAccessor
}

func NewQueryAccessor(query url.Values, secretKey string, setters ...Setter) (*QueryAccessor, error) {
	qa := &QueryAccessor{
		baseAccessor: newBaseAccessor(),
	}
	for key, vs := range query {
		v := vs[0]
		if len(v) == 0 {
			return nil, fmt.Errorf("%w: %s", ErrArgLack, key)
		}
		qa.args.append(key, v)
	}
	qa.args.append(secretKeyTag, secretKey)
	if len(qa.args.kv[nonceTag]) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrArgLack, nonceTag)
	}
	if len(qa.args.kv[timestampTag]) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrArgLack, timestampTag)
	}
	if len(qa.args.kv[signatureTag]) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrArgLack, signatureTag)
	}

	for _, setter := range setters {
		if err := setter(&qa.baseAccessor); err != nil {
			return nil, err
		}
	}

	return qa, nil
}
