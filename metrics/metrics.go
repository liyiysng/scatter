// Package metrics 使用 prometheus 检测指标
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var userInput = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "user_input",
	}, []string{"method", "endpoint"},
)

func init() {
	prometheus.MustRegister(userInput)
}

// StartListen 开始启动指标监听
func StartListen(addr string) error {

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}
	return nil
}

// UserInput 用户输入
func UserInput(method, endpoint string) {
	userInput.With(prometheus.Labels{"method": method, "endpoint": endpoint}).Inc()
}
