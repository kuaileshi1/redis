// @Title 连接状态检测
// @Description
// @Author shigx 2021/6/17 6:12 下午
package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"net"
	"sync/atomic"
)

const (
	ActionEject = iota
	ActionRecover
)

type StateClient struct {
	*redis.Client
	id             string
	actionListener chan<- Action
	evicted        bool
	failureCount   int32
}

type Action struct {
	id     string
	action int
}

func newStateClient(opt *redis.Options, actionListener chan<- Action) *StateClient {
	cli := redis.NewClient(opt)
	s := &StateClient{
		id:             opt.Addr,
		Client:         cli,
		failureCount:   0,
		actionListener: actionListener,
	}
	cli.AddHook(newFailureHook(s))

	return s
}

func (s *StateClient) onFailure() {
	if s.evicted {
		return
	}

	failureCount := atomic.AddInt32(&s.failureCount, 1)
	if failureCount > failureLimit {
		s.evicted = true
		s.actionListener <- Action{
			id:     s.id,
			action: ActionEject,
		}
	}
}

func (s *StateClient) onSuccess() {
	if !s.evicted {
		return
	}
	atomic.StoreInt32(&s.failureCount, 0)
	s.evicted = false
	s.actionListener <- Action{
		id:     s.id,
		action: ActionRecover,
	}
}

func (s *StateClient) ping() bool {
	result := s.Ping(context.Background())
	if result == nil {
		return false
	}

	pong, err := result.Result()
	if err != nil {
		return false
	}
	return pong == "PONG"
}

type HealthHook struct {
	*StateClient
}

func newFailureHook(s *StateClient) *HealthHook {
	return &HealthHook{
		StateClient: s,
	}
}

func (hook *HealthHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (hook *HealthHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if hook.isNetworkError(cmd.Err()) {
		hook.onFailure()
	} else {
		hook.onSuccess()
	}

	return nil
}

func (hook *HealthHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (hook *HealthHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		if hook.isNetworkError(cmd.Err()) {
			hook.StateClient.onFailure()
			return nil
		}
	}
	hook.onSuccess()
	return nil
}

func (hook *HealthHook) isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(net.Error)

	return ok
}
