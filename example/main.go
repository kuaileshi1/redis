// @Title 调用演示
// @Description 简单使用实例
// @Author shigx 2021/6/18 5:09 下午
package main

import (
	"context"
	"fmt"
	"github.com/kuaileshi1/redis"
	"time"
)

func main() {
	// 1、主从模式
	cfg := redis.ConfigRedis{
		Driver:           redis.DriverReplication,
		Master:           "127.0.0.1:6379",
		Slaves:           "127.0.0.1:6379,127.0.0.1:6379",
		Password:         "",
		ReadonlyPassword: "",
		Options:          nil,
	}

	engine, _ := redis.NewClient(&cfg)
	// 设置缓存
	engine.Set(context.Background(), "key1", "2222", time.Hour)
	// 读取缓存
	if result, err := engine.Get(context.Background(), "key1").Result(); err != redis.Nil {
		fmt.Println("Slave:", result)
	}
	// 指定主库读取缓存
	if result, err := engine.Master().Get(context.Background(), "key1").Result(); err != redis.Nil {
		fmt.Println("Master:", result)
	}

	// 2、集群模式
	cfg1 := redis.ConfigRedis{
		Driver: redis.DriverCluster,
		Master: "127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003,127.0.0.1:7000,127.0.0.1:7004,127.0.0.1:7005",
	}

	engine1, _ := redis.NewClient(&cfg1)
	// 设置缓存
	engine1.Set(context.Background(), "key2", "3333", time.Hour)
	// 读取缓存
	if result, err := engine1.Get(context.Background(), "key2").Result(); err != redis.Nil {
		fmt.Println("get:", result)
	}
}
