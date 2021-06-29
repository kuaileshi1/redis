// @Title 主从复制驱动处理
// @Description 实例化对象操作
// @Author shigx 2021/6/17 5:53 下午
package redis

import (
	"github.com/go-redis/redis/v8"
	"strings"
)

type Replication struct {
	master *redis.Client
	slaves *MultiClient
}

// 实例化
func NewReplication(cfg *ConfigRedis) *Replication {
	if cfg.Options == nil {
		cfg.Options = &redis.Options{}
	}
	addr := cfg.Master
	// 配置支持","分割，主从模式下如果配置多个的话只取第一个
	if addrs := strings.Split(addr, ","); len(addrs) > 1 {
		addr = addrs[0]
	}
	cfg.Options.Addr = addr
	cfg.Options.Password = cfg.Password

	master := redis.NewClient(cfg.Options)
	slaves := NewMultiClient(cfg)

	return &Replication{
		master: master,
		slaves: slaves,
	}
}

// 返回主客户端
func (r *Replication) Master() Client {
	return r.master
}

func (r *Replication) Slave() Client {
	return r.slaves.GetConn()
}

// 添加执行钩子
func (r *Replication) AddHook(hook redis.Hook) {
	r.master.AddHook(hook)
	r.slaves.AddHook(hook)
}
