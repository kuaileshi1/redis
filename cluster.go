// @Title 集群驱动处理
// @Description 实例化集群客户端
// @Author shigx 2021/6/17 5:19 下午
package redis

import (
	"github.com/go-redis/redis/v8"
	"strings"
)

type ClusterSplit struct {
	client *redis.ClusterClient
}

// 实例化连接
func NewCluster(cfg *ConfigRedis) *ClusterSplit {
	addr := strings.Split(cfg.Master, ",")
	if cfg.Options == nil {
		cfg.Options = &redis.Options{}
	}

	clusterOptions := &redis.ClusterOptions{
		Addrs:              addr,
		Password:           cfg.Password,
		DialTimeout:        cfg.Options.DialTimeout,
		ReadTimeout:        cfg.Options.ReadTimeout,
		WriteTimeout:       cfg.Options.WriteTimeout,
		PoolSize:           cfg.Options.PoolSize,
		MinIdleConns:       cfg.Options.MinIdleConns,
		MaxConnAge:         cfg.Options.MaxConnAge,
		PoolTimeout:        cfg.Options.PoolTimeout,
		IdleTimeout:        cfg.Options.IdleTimeout,
		IdleCheckFrequency: cfg.Options.IdleCheckFrequency,
	}

	client := redis.NewClusterClient(clusterOptions)

	return &ClusterSplit{
		client: client,
	}
}

// 集群模式不分主从，直接返回连接
func (c *ClusterSplit) Master() Client {
	return c.client
}

// 集群模式不分主从，直接返回连接
func (c *ClusterSplit) Slave() Client {
	return c.client
}

// 设置执行钩子
func (c *ClusterSplit) AddHook(hook redis.Hook) {
	c.client.AddHook(hook)
}
