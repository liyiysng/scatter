// Package config 配置信息
package config

import (
	"fmt"
	"strings"

	"github.com/liyiysng/scatter/logger"
	"github.com/spf13/viper"
)

var (
	myLog = logger.Component("config")
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

// ConfigValid 配置是否正确
func (c *Config) ConfigValid() error {

	if c.Get("scatter.es.write_index") != c.Get("scatter.es.template.settings.index.lifecycle.rollover_alias") {
		return fmt.Errorf("scatter.es.write_index must same as scatter.es.template.settings.index.lifecycle.rollover_alias")
	}

	return nil
}

func (c *Config) fillDefaultValues() {

	writeIndexName := "scatter_write"

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

		// es settings ////////////////////////////////////////////////////////////////////////////////////////////////////////
		////////////////////////////////////////////////////////////////
		"scatter.es.url":             []string{"http://localhost:9200"},
		"scatter.es.init_index_name": "<scatter-{now/d}-1>",
		"scatter.es.write_index":     writeIndexName,
		// 生命周期
		"scatter.es.lifecycle_name": "scatter_lifecycle_policy",
		// hot
		"scatter.es.lifecycle.policy.phases.hot.min_age":                       "0ms",
		"scatter.es.lifecycle.policy.phases.hot.actions.rollover.max_age":      "30d",  // 热数据时间
		"scatter.es.lifecycle.policy.phases.hot.actions.rollover.max_size":     "50gb", // 热数据大小
		"scatter.es.lifecycle.policy.phases.hot.actions.rollover.max_docs":     10000,  // 热数据个数
		"scatter.es.lifecycle.policy.phases.hot.actions.set_priority.priority": 100,
		// delete
		"scatter.es.lifecycle.policy.phases.delete.min_age":        "100d",     // 多久后删除
		"scatter.es.lifecycle.policy.phases.delete.actions.delete": struct{}{}, // 若要删除,徐保留
		////////////////////////////////////////////////////////////////
		// 模板
		"scatter.es.template_name":           "scatter_template",
		"scatter.es.template.index_patterns": "scatter-*",
		// 设置
		"scatter.es.template.settings.number_of_shards":               1,
		"scatter.es.template.settings.number_of_replicas":             1,
		"scatter.es.template.settings.index.lifecycle.name":           "scatter_lifecycle_policy",
		"scatter.es.template.settings.index.lifecycle.rollover_alias": writeIndexName,
		// 别名(别名为scatters)
		"scatter.es.template.aliases.scatters": struct{}{},
		// mappings

		// register settings//////////////////////////////////////////////////////////////////////////////////////////////
		"scatter.register.grpc_check_interval": "5s",     // 健康检查
		"scatter.register.default":             "consul", //默认服务注册/发现
		// consul settings
		"scatter.register.consul.addrs":   []string{"127.0.0.1:8500"},
		"scatter.register.consul.timeout": "0",

		// service setting////////////////////////////////////////////////////////////////////////////////////////////////
		"scatter.service.srvints": &Service{
			SelectPolicy: "consistent_hash",
			Meta:         map[string]string{"foo": "bar"},
		},
		"scatter.service.srvstrings.selectpolicy": "consistent_hash",
		"scatter.service.srvstrings.meta":         map[string]string{"foo": "bar"},
		// sub service setting//////////////////////////////////////////////////////////////////////////////////////////
		"scatter.subservice.foo.selectpolicy": "consistent_hash",
		"scatter.subservice.foo1": &SubService{
			SelectPolicy: "consistent_hash",
		},
		// grpc node settings///////////////////////////////////////////////////////////////////////////////////////////
		"scatter.gnode.dial.timeout": "10s", // 链接超时
	}

	for param := range defaultsMap {
		if c.Get(param) == nil {
			c.SetDefault(param, defaultsMap[param])
		}
	}
}
