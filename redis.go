// @Title 客户端实现
// @Description 实现接口方法
// @Author shigx 2021/6/17 4:18 下午
package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var redisPool sync.Map

// 连接放入连接池
func PutClient(namespace string, client *splitClient) {
	redisPool.Store(namespace, client)
}

// 从连接池取出连接
func GetClient(namespace string) (*splitClient, error) {
	v, ok := redisPool.Load(namespace)
	if !ok {
		return nil, errors.New(namespace + " redis未定义")
	}
	return v.(*splitClient), nil
}

// 实例化客户端
func NewClient(cfg *ConfigRedis) (*splitClient, error) {
	if cfg.Driver == DriverReplication {
		return &splitClient{
			splitter: NewReplication(cfg),
		}, nil
	}

	if cfg.Driver == DriverCluster {
		return &splitClient{
			splitter: NewCluster(cfg),
		}, nil
	}

	return nil, errors.New("unSupported driver:" + cfg.Driver)
}

type splitClient struct {
	splitter Splitter
}

// 获取主库连接
func (s *splitClient) Master() Client {
	return s.splitter.Master()
}

// 返回从库连接
func (s *splitClient) Slave() Client {
	return s.splitter.Slave()
}

// Get
func (s *splitClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}

	return cli.Get(ctx, key)
}

func (s *splitClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}

	return cli.StrLen(ctx, key)
}

func (s *splitClient) MGet(ctx context.Context, key ...string) *redis.SliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.MGet(ctx, key...)
}

func (s *splitClient) Set(ctx context.Context, key string, value interface{}, seconds time.Duration) *redis.StatusCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.Set(ctx, key, value, seconds)
}

func (s *splitClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.Decr(ctx, key)
}

func (s *splitClient) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.DecrBy(ctx, key, decrement)
}

func (s *splitClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.Incr(ctx, key)
}

func (s *splitClient) IncrBy(ctx context.Context, key string, increment int64) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.IncrBy(ctx, key, increment)
}

func (s *splitClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HExists(ctx, key, field)
}

func (s *splitClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HGet(ctx, key, field)
}

func (s *splitClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HGetAll(ctx, key)
}

func (s *splitClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HKeys(ctx, key)
}

func (s *splitClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HLen(ctx, key)
}

func (s *splitClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HMGet(ctx, key, fields...)
}

func (s *splitClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.HVals(ctx, key)
}

func (s *splitClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}

	return cli.HDel(ctx, key, fields...)
}

func (s *splitClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.HIncrBy(ctx, key, field, incr)
}

func (s *splitClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.HMSet(ctx, key, values...)
}

func (s *splitClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.HSet(ctx, key, values...)
}

func (s *splitClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.Del(ctx, keys...)
}

func (s *splitClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.Exists(ctx, keys...)
}

func (s *splitClient) Expire(ctx context.Context, key string, seconds time.Duration) *redis.BoolCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.Expire(ctx, key, seconds)
}

func (s *splitClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.Type(ctx, key)
}

func (s *splitClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.TTL(ctx, key)
}

func (s *splitClient) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LIndex(ctx, key, index)
}

func (s *splitClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.LLen(ctx, key)
}

func (s *splitClient) LRange(ctx context.Context, key string, start, end int64) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.LRange(ctx, key, start, end)
}

func (s *splitClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LInsert(ctx, key, op, pivot, value)
}

func (s *splitClient) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LInsertBefore(ctx, key, pivot, value)
}

func (s *splitClient) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LInsertAfter(ctx, key, pivot, value)
}

func (s *splitClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LPop(ctx, key)
}

func (s *splitClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LPush(ctx, key, values...)
}

func (s *splitClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LRem(ctx, key, count, value)
}

func (s *splitClient) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LSet(ctx, key, index, value)
}

func (s *splitClient) LTrim(ctx context.Context, key string, start, end int64) *redis.StatusCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.LTrim(ctx, key, start, end)
}

func (s *splitClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.RPop(ctx, key)
}

func (s *splitClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.RPush(ctx, key, values...)
}

func (s *splitClient) ZRange(ctx context.Context, key string, start, end int64) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRange(ctx, key, start, end)
}

func (s *splitClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.ZAdd(ctx, key, members...)
}

func (s *splitClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZCard(ctx, key)
}

func (s *splitClient) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZCount(ctx, key, min, max)
}

func (s *splitClient) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRangeByLex(ctx, key, opt)
}

func (s *splitClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRangeByScore(ctx, key, opt)
}

func (s *splitClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRangeByScoreWithScores(ctx, key, opt)
}

func (s *splitClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.ZIncrBy(ctx, key, increment, member)
}

func (s *splitClient) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRank(ctx, key, member)
}

func (s *splitClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.ZRem(ctx, key, members...)
}

func (s *splitClient) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.ZRemRangeByScore(ctx, key, min, max)
}

func (s *splitClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRemRangeByRank(ctx, key, start, stop)
}

func (s *splitClient) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.ZRemRangeByLex(ctx, key, min, max)
}

func (s *splitClient) ZRevRange(ctx context.Context, key string, start, end int64) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRevRange(ctx, key, start, end)
}

func (s *splitClient) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRevRangeByScore(ctx, key, opt)
}

func (s *splitClient) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (s *splitClient) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZRevRank(ctx, key, member)
}

func (s *splitClient) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZScore(ctx, key, member)
}

func (s *splitClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.ZScan(ctx, key, cursor, match, count)
}

func (s *splitClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.SMembers(ctx, key)
}

func (s *splitClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.SAdd(ctx, key, members...)
}

func (s *splitClient) AddHook(hook redis.Hook) {
	switch s.splitter.(type) {
	case *ClusterSplit:
		s.splitter.(*ClusterSplit).AddHook(hook)
	case *Replication:
		s.splitter.(*Replication).AddHook(hook)
	}
}

func (s *splitClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.SetNX(ctx, key, value, expiration)
}

func (s *splitClient) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.MSetNX(ctx, values...)
}

func (s *splitClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.GetBit(ctx, key, offset)
}
func (s *splitClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.SetBit(ctx, key, offset, value)
}
func (s *splitClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.BitCount(ctx, key, bitCount)
}
func (s *splitClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.BitOpAnd(ctx, destKey, keys...)
}
func (s *splitClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.BitOpOr(ctx, destKey, keys...)
}
func (s *splitClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.BitOpXor(ctx, destKey, keys...)
}
func (s *splitClient) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	cli := s.splitter.Master()
	if cli == nil {
		return nil
	}
	return cli.BitOpNot(ctx, destKey, key)
}
func (s *splitClient) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *redis.IntCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.BitPos(ctx, key, bit, pos...)
}
func (s *splitClient) BitField(ctx context.Context, key string, args ...interface{}) *redis.IntSliceCmd {
	cli := s.splitter.Slave()
	if cli == nil {
		return nil
	}
	return cli.BitField(ctx, key, args...)
}
