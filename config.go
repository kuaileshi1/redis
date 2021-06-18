// @Title 配置文件
// @Description 配置定义
// @Author shigx 2021/6/17 3:39 下午
package redis

import (
	"github.com/go-redis/redis/v8"
	"time"
)

// 驱动类型定义
const (
	DriverCluster     = "cluster"
	DriverReplication = "replication"
)

type ConfigRedis struct {
	Driver             string         `yaml:"driver"`           // replication:主从复制；cluster:集群
	Master             string         `yaml:"master"`           // 主节点IP和端口,多个用","分割，eg: 127.0.0.1:6379
	Slaves             string         `yaml:"slaves"`           // 从节点IP和端口,多个用","分割，eg: "127.0.0.1:6379,127.0.0.1:6380,127.0.0.1:6381"
	Password           string         `yaml:"password"`         // 主节点授权密码
	ReadonlyPassword   string         `yaml:"readonlyPassword"` // 从节点授权密码
	Options            *redis.Options // redis连接池等配置信息
	ServerFailureLimit int            // 服务失败次数限制，达到该限制剔除
	ServerRetryTimeout time.Duration  // 服务失败重试超时时间
}
