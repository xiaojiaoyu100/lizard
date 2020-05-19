package tokenbucket

import "github.com/go-redis/redis/v7"

var (
	_ rediser = (*redis.Client)(nil)
	_ rediser = (*redis.Ring)(nil)
	_ rediser = (*redis.ClusterClient)(nil)
)

type rediser interface {
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(hashes ...string) *redis.BoolSliceCmd
	ScriptLoad(script string) *redis.StringCmd
}
