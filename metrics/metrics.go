// Package metrics 使用 prometheus 检测指标
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartListen 开始启动指标监听
func StartListen(addr string) error {

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}
	return nil
}
