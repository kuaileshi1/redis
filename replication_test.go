// @Title 主从模式
// @Description 测试用例
// @Author shigx 2021/6/20 12:35 上午
package redis

import (
	"context"
	"log"
	"testing"
)

func TestReplication(t *testing.T) {
	cli := NewReplication(&ConfigRedis{
		Driver: DriverReplication,
		Master: "127.0.0.1:6379",
		Slaves: "127.0.0.1:6379,127.0.0.1:6379",
	})

	cli.Master().Set(context.Background(), "key1", "v1", 0)
	// 考虑所有的slave 全部挂掉的情况
	client := cli.Slave()
	if client == nil {
		return
	}
	cmd := client.Get(context.Background(), "key1")
	if cmd != nil {
		log.Println(cmd.Result())
	}
}
