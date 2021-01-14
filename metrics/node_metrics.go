package metrics

import (
	"sync"

	"github.com/liyiysng/scatter/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	myLog = logger.Component("metrics")
)

type nodeMetrics struct {
	sync.Mutex
	connCount map[string] /*id*/ prometheus.Gauge
}

func (nm *nodeMetrics) GetConnCount(id string) prometheus.Gauge {
	nm.Lock()
	defer nm.Unlock()
	if g, ok := nm.connCount[id]; ok {
		return g
	}
	g := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "node",
			Name:        "conn_count",
			Help:        "node 当前链接数量",
			ConstLabels: prometheus.Labels{"id": id},
		},
	)
	nm.connCount[id] = g
	err := prometheus.Register(g)
	if err != nil {
		myLog.Error(err)
	}
	return g
}

func newNodeMetrics() *nodeMetrics {
	return &nodeMetrics{
		connCount: make(map[string]prometheus.Gauge),
	}
}

var nm = newNodeMetrics()

// ReportNodeNewCommingConn 新的链接
func ReportNodeNewCommingConn(id string) {
	nm.GetConnCount(id).Inc()
}

// ReportNodeDisConn 链接断开
func ReportNodeDisConn(id string) {
	nm.GetConnCount(id).Dec()
}
