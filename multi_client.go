// @Title 多实例处理
// @Description 实例化多实例
// @Author shigx 2021/6/17 6:08 下午
package redis

import (
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

const (
	failureLimit = 1
)

type MultiClient struct {
	avails      map[string]interface{}
	clients     map[string]*StateClient
	healthCheck time.Duration
	stateCh     chan Action
	mu          sync.RWMutex
	stopCh      chan struct{}
}

func NewMultiClient(cfg *ConfigRedis) *MultiClient {
	slavePassword := cfg.Password
	if cfg.ReadonlyPassword != "" {
		slavePassword = cfg.Password
	}
	var slaves []string
	if cfg.Slaves == "" {
		slaves = strings.Split(cfg.Master, ",")
	} else {
		slaves = strings.Split(cfg.Slaves, ",")
	}
	if cfg.ServerFailureLimit == 0 {
		cfg.ServerFailureLimit = failureLimit
	}
	if cfg.ServerRetryTimeout == 0 {
		cfg.ServerRetryTimeout = time.Second
	}

	pool := &MultiClient{
		stopCh:      make(chan struct{}, 0),
		stateCh:     make(chan Action),
		clients:     make(map[string]*StateClient, len(slaves)),
		avails:      make(map[string]interface{}, len(slaves)),
		healthCheck: cfg.ServerRetryTimeout,
	}
	options := cfg.Options
	for _, slave := range slaves {
		slaveOption := *options
		slaveOption.Addr = slave
		slaveOption.Password = slavePassword
		stateClient := newStateClient(&slaveOption, pool.stateCh)
		pool.clients[slave] = stateClient
		pool.avails[slave] = nil
	}

	go pool.listenStateChange()
	go pool.recoveryCheck()

	return pool
}

func (m *MultiClient) AddHook(hook redis.Hook) {
	for _, c := range m.clients {
		c.AddHook(hook)
	}
}

func (m *MultiClient) listenStateChange() {
	for action := range m.stateCh {
		m.changeState(action)
	}
}

func (m *MultiClient) changeState(action Action) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch action.action {
	case ActionEject:
		delete(m.avails, action.id)
	case ActionRecover:
		m.avails[action.id] = nil
	}
}

func (m *MultiClient) recover() {
	for _, slave := range m.clients {
		slave.ping()
	}
}

func (m *MultiClient) recoveryCheck() {
	ticker := time.NewTicker(m.healthCheck)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.recover()
		}
	}
}

func (m *MultiClient) close() {
	close(m.stopCh)
	for _, cli := range m.clients {
		cli.Client.Close()
	}
}

func (m *MultiClient) GetConn() *StateClient {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.avails) == 0 {
		return nil
	}
	// 这里有点偷懒，借助map循环的无序性返回循环的第一个值。具体每个key出现的概率不能保证
	for id := range m.avails {
		return m.clients[id]
	}

	return nil
}
