package tokenbucket

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/xiaojiaoyu100/lizard/timekit"
)

const script = `
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local tokenNum = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local num = tonumber(ARGV[4])
local expiration = ARGV[5]
local obj = {
tn=tokenNum,
ts=now
}

local value = redis.call("get", key)
if value then
  obj = cjson.decode(value)
end

local incr = math.floor((now - obj.ts) / rate)
if incr > 0 then
  obj.tn = math.min(obj.tn + incr, tokenNum)
  obj.ts = obj.ts + incr * rate
end

if obj.tn >= num then
  obj.tn = obj.tn - num
  obj.ts = string.format("%.f", obj.ts)
  if redis.call("set", key, cjson.encode(obj), "EX", expiration) then
    return 1
  end
end

return 0
`

var scriptDigest string

func init() {
	s := sha1.New()
	io.WriteString(s, script)
	scriptDigest = hex.EncodeToString(s.Sum(nil))
}

// TokenBucket stands for a token bucket.
type TokenBucket struct {
	client     *redis.Client // redis client
	Key        string        // redis key
	TokenNum   int64         // token bucket size
	Rate       time.Duration // the rate of putting token into bucket
	Expiration int64         // redis key expiration in seconds
}

// New returns an instance of TokenBucket
func New(client *redis.Client, key string, tokenNum int64, rate time.Duration, expiration int64) (*TokenBucket, error) {
	h := sha1.New()
	_, err := io.WriteString(h, script)
	if err != nil {
		return nil, err
	}

	if timekit.DurationToMillis(rate) == 0 {
		return nil, errors.New("wrong rate")
	}

	return &TokenBucket{
		client:     client,
		Key:        key,
		TokenNum:   tokenNum,
		Rate:       rate,
		Expiration: expiration,
	}, nil
}

func (tb *TokenBucket) eva(script string, key string, argv ...interface{}) (int64, error) {
	ret, err := tb.client.Eval(script, []string{key}, argv...).Result()
	if err != nil {
		return 0, err
	}
	return ret.(int64), nil
}

func (tb *TokenBucket) evaSha1(sha1 string, key string, argv ...interface{}) (int64, error) {
	ret, err := tb.client.EvalSha(sha1, []string{key}, argv...).Result()
	if err != nil {
		return 0, err
	}
	return ret.(int64), nil
}

// Consume consumes the number of token in the token bucket.
func (tb *TokenBucket) Consume(num int64) (bool, error) {
	if num > tb.TokenNum {
		return false, errors.New("token is not enough")
	}
	ok, err := tb.evaSha1(scriptDigest, tb.Key, timekit.DurationToMillis(tb.Rate), tb.TokenNum, timekit.NowInMillis(), num, tb.Expiration)
	// NOSCRIPT 这个error是稳定的 see https://redis.io/commands/eval
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT") {
		ok, err := tb.eva(script, tb.Key, timekit.DurationToMillis(tb.Rate), tb.TokenNum, timekit.NowInMillis(), num, tb.Expiration)
		if err != nil {
			return false, err
		}
		return ok == 1, nil
	}
	if err != nil {
		return false, err
	}
	return ok == 1, nil
}
