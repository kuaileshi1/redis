// @Title redis使用
// @Description 测试用例
// @Author shigx 2021/6/20 12:17 上午
package redis

import (
	"context"
	"testing"
	"time"
)

func newClient() (*splitClient, error) {
	return NewClient(&ConfigRedis{
		Driver:           DriverReplication,
		Master:           "127.0.0.1:6379",
		Slaves:           "127.0.0.1:6379,127.0.0.1:6379",
		Password:         "",
		ReadonlyPassword: "",
	})
}

func TestSplitClient(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Log(err)
		return
	}
	cli.Set(context.Background(), "key1", "222222", time.Hour)
	cmd := cli.Get(context.Background(), "key1")
	if cmd != nil {
		t.Log(cmd.Result())
	}
}

func TestMaster(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Log(err)
		return
	}

	cli.Master().Set(context.Background(), "key2", "333333", 0)
	cmd := cli.Master().Get(context.Background(), "key2")
	if cmd != nil {
		t.Log(cmd.Result())
	}
}

func TestSlave(t *testing.T) {
	cli, err := newClient()
	if err != nil {
		t.Log(err)
		return
	}

	setResult := cli.Slave().Set(context.Background(), "key3", "44444", 0)
	_, err = setResult.Result()
	// execute set on readonly node error
	if err != nil {
		t.Log(err)
	}
	cmd := cli.Slave().Get(context.Background(), "key3")
	if cmd != nil {
		t.Log(cmd.Result())
	}
}
