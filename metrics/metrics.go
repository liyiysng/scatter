// Package metrics 使用 prometheus 检测指标
package metrics

import (
	"runtime"
	"time"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("metrics")
)

var (
	// ResponseTime reports the response time of handlers and rpc
	ResponseTime = "response_time_ns"
	// ConnectedClients represents the number of current connected clients in frontend servers
	ConnectedClients = "connected_clients"
	// ProcessDelay reports the message processing delay to handle the messages at the handler service
	ProcessDelay = "handler_delay_ns"
	// Goroutines reports the number of goroutines
	Goroutines = "goroutines"
	// HeapSize reports the size of heap
	HeapSize = "heapsize"
	// HeapObjects reports the number of allocated heap objects
	HeapObjects = "heapobjects"
)

// Reporter interface
type Reporter interface {
	ReportCount(metric string, tags map[string]string, count float64) error
	ReportSummary(metric string, tags map[string]string, value float64) error
	ReportGauge(metric string, tags map[string]string, value float64) error
}

// ReportSysMetrics reports sys metrics
func ReportSysMetrics(reporters []Reporter, period time.Duration) {
	for {
		for _, r := range reporters {
			num := runtime.NumGoroutine()
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)

			myLog.Info("ReportSysMetrics")
			err := r.ReportGauge(Goroutines, map[string]string{}, float64(num))
			if err != nil {
				myLog.Error(err)
				return
			}
			r.ReportGauge(HeapSize, map[string]string{}, float64(m.Alloc))
			if err != nil {
				myLog.Error(err)
				return
			}
			r.ReportGauge(HeapObjects, map[string]string{}, float64(m.HeapObjects))
			if err != nil {
				myLog.Error(err)
				return
			}
		}

		time.Sleep(period)
	}
}
