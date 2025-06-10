package config

import "sync"

type Config struct {
	ShowReasoning bool
	mu            sync.RWMutex
}

var globalConfig = &Config{}

func SetShowReasoning(v bool) {
	globalConfig.mu.Lock()
	defer globalConfig.mu.Unlock()
	globalConfig.ShowReasoning = v
}

func GetShowReasoning() bool {
	globalConfig.mu.RLock()
	defer globalConfig.mu.RUnlock()
	return globalConfig.ShowReasoning
}
