// @Title 集群模式
// @Description 集群模式使用
// @Author shigx 2021/6/19 11:33 下午
package redis

import (
	"context"
	"testing"
	"time"
)

func TestCluster(t *testing.T) {
	cli := NewCluster(&ConfigRedis{
		Driver: DriverCluster,
		Master: "127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003,127.0.0.1:7000,127.0.0.1:7004,127.0.0.1:7005",
	})

	cli.Master().Set(context.Background(), "key1", "11111111", time.Hour)

	result, err := cli.Slave().Get(context.Background(), "key1").Result()
	if err == nil {
		t.Log(result)
	}
}
