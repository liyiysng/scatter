package config

import "sync"

var defaultConfig *Config

var initConfigOnce sync.Once

// GetConfig 获取单例配置
func GetConfig() *Config {

	initConfigOnce.Do(func() {
		defaultConfig = NewConfig()
	})

	return defaultConfig
}
