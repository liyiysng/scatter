// Package config 配置信息
package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config is a wrapper around a viper config
type Config struct {
	*viper.Viper
}

// NewConfig creates a new config with a given viper config if given
func NewConfig(cfgs ...*viper.Viper) *Config {
	var cfg *viper.Viper
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	} else {
		cfg = viper.New()
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	c := &Config{Viper: cfg}
	c.fillDefaultValues()
	return c
}

func (c *Config) fillDefaultValues() {
	defaultsMap := map[string]interface{}{
		"scatter.game": "scatter",
		// metrics config
		"scatter.metrics.prometheus.collect_type": "push", // push , listen
		// metrics pusher
		"scatter.metrics.prometheus.pusher.addr":    ":9091",
		"scatter.metrics.prometheus.pusher.inteval": "1s",
		// metrics listen
		"scatter.metrics.prometheus.port": 8888,
		"scatter.metrics.additionalTags":  map[string]string{},
		"scatter.metrics.constTags":       map[string]string{},
		"scatter.metrics.custom":          map[string]interface{}{},
	}

	for param := range defaultsMap {
		if c.Get(param) == nil {
			c.SetDefault(param, defaultsMap[param])
		}
	}
}
