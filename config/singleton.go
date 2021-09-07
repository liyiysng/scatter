package config

import (
	"os"
	"sync"

	"github.com/liyiysng/scatter/logger"
	"github.com/spf13/viper"
)

var defaultConfig *Config

var initConfigOnce sync.Once

// GetConfig 获取单例配置
func GetConfig() *Config {

	initConfigOnce.Do(func() {
		configPath := os.Getenv("SCATTER_CONFIG_PATH")
		if len(configPath) == 0 {
			defaultConfig = NewConfig()
		} else {
			v := viper.New()
			v.SetConfigFile(configPath)
			if myLog.V(logger.VIMPORTENT) {
				myLog.Infof("use config %s ", configPath)
			}
			v.ReadInConfig()
			defaultConfig = NewConfig(v)
		}

	})

	return defaultConfig
}
