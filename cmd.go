// @Title 接口定义
// @Description redis相关接口定义
// @Author shigx 2021/6/17 4:02 下午
package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type SplitClient interface {
	Splitter
	Client
}

type Splitter interface {
	Master() Client
	Slave() Client
}

// 客户端支持方法接口定义
type Client interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	StrLen(ctx context.Context, key string) *redis.IntCmd
	MGet(ctx context.Context, key ...string) *redis.SliceCmd
	Set(ctx context.Context, key string, value interface{}, seconds time.Duration) *redis.StatusCmd
	Decr(ctx context.Context, key string) *redis.IntCmd
	DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	IncrBy(ctx context.Context, key string, increment int64) *redis.IntCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd
	HKeys(ctx context.Context, key string) *redis.StringSliceCmd
	HLen(ctx context.Context, key string) *redis.IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd
	HVals(ctx context.Context, key string) *redis.StringSliceCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd
	HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Expire(ctx context.Context, key string, seconds time.Duration) *redis.BoolCmd
	Type(ctx context.Context, key string) *redis.StatusCmd
	TTL(ctx context.Context, key string) *redis.DurationCmd
	LIndex(ctx context.Context, key string, index int64) *redis.StringCmd
	LLen(ctx context.Context, key string) *redis.IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd
	LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd
	LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd
	LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd
	LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZCard(ctx context.Context, key string) *redis.IntCmd
	ZCount(ctx context.Context, key, min, max string) *redis.IntCmd
	ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd
	ZRank(ctx context.Context, key, member string) *redis.IntCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd
	ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd

	ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd
	ZRevRank(ctx context.Context, key, member string) *redis.IntCmd
	ZScore(ctx context.Context, key, member string) *redis.FloatCmd
	ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd
	SMembers(ctx context.Context, key string) *redis.StringSliceCmd
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd
	AddHook(hook redis.Hook)

	GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd
	SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd
	BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd
	BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd
	BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd
	BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd
	BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd
	BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd
	BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd
}
