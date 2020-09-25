package apiaccessor

import "github.com/go-redis/redis/v7"

type Setter func(b *baseAccessor) error

func WithEvalSignatureFunc(e EvalSignature) Setter {
	return func(b *baseAccessor) error {
		b.evalSignatureFunc = e
		return nil
	}
}

var checkCountScript = redis.NewScript(`
local key = KEYS[1]
local threshold_count = tonumber(ARGV[1])
local threshold_sec = tonumber(ARGV[2])

local a_sec_millisecond = 1000
local over = 1
local no = -1

local count = redis.call("INCR", key)
if tonumber(count) == 1 then
	redis.call("PEXPIRE", key, threshold_sec * a_sec_millisecond)
end

if tonumber(count) > threshold_count then
	return over
end
return no
`)

type KeyGen func(nonce string) (key string)

func WithGeneralRedisNonceChecker(client redis.Cmdable, sec int64, keyGenFunc KeyGen) Setter {
	return func(b *baseAccessor) error {
		b.nonceChecker = func(nonce string) error {
			key := keyGenFunc(b.args.kv[nonceTag])
			re, err := checkCountScript.Run(client, []string{key}, 1, sec).Int()
			if err != nil {
				return err
			}
			if re == 1 {
				return ErrNonceUsed
			}
			return nil
		}
		return nil
	}
}

func WithNonceChecker(nc NonceChecker) Setter {
	return func(b *baseAccessor) error {
		b.nonceChecker = nc
		return nil
	}
}
