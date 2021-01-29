package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	_readBytesCountTotal  = "read_bytes_count_total"
	_writeBytesCountTotal = "write_bytes_count_total"
)

// CounterMeta 计数指标
// key 作为lable的值
type CounterMeta struct {
	SeverName string
	NID       string
	Name      string
	Help      string
	Lables    []string
}

type nodeMetrics struct {
	sync.Mutex
	connCount map[string] /*nid*/ prometheus.Gauge

	counter map[string] /*counter name*/ map[interface{}]prometheus.Counter

	counterMeta map[string]*CounterMeta
}

func (nm *nodeMetrics) addCounterMeta(cmeta *CounterMeta) {
	nm.counterMeta[cmeta.Name] = cmeta
}

func (nm *nodeMetrics) getCounter(nid string, name string, key interface{}) prometheus.Counter {
	nm.Lock()
	defer nm.Unlock()
	if cs, ok := nm.counter[name]; ok {
		if c, cok := cs[key]; cok {
			return c
		}
	} else {
		nm.counter[name] = make(map[interface{}]prometheus.Counter)
	}
	c := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "scatter",
			Subsystem:   nid,
			Name:        name,
			Help:        "当前读字节数",
			ConstLabels: prometheus.Labels{"sid": fmt.Sprintf("%v", key)},
		},
	)
	err := prometheus.Register(c)
	if err != nil {
		myLog.Error(err)
	}
	nm.counter[name][key] = c
	return c
}

func (nm *nodeMetrics) getConnCount(nid string) prometheus.Gauge {
	nm.Lock()
	defer nm.Unlock()
	if g, ok := nm.connCount[nid]; ok {
		return g
	}

	g := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "node",
			Name:        "conn_count",
			Help:        "node 当前链接数量",
			ConstLabels: prometheus.Labels{"id": nid},
		},
	)
	nm.connCount[nid] = g
	err := prometheus.Register(g)
	if err != nil {
		myLog.Error(err)
	}
	return g
}

func newNodeMetrics() *nodeMetrics {
	return &nodeMetrics{
		connCount: make(map[string]prometheus.Gauge),
		counter:   make(map[string] /*counter name*/ map[interface{}]prometheus.Counter),
	}
}

var nm = newNodeMetrics()

// ReportNodeNewCommingConn 新的链接
func ReportNodeNewCommingConn(nid string) {
	nm.getConnCount(nid).Inc()
}

// ReportNodeDisConn 链接断开
func ReportNodeDisConn(nid string) {
	nm.getConnCount(nid).Dec()
}

// ReportConnReadBytesCount 链接读取字节数
func ReportConnReadBytesCount(nid string, sid int64, rBytes int) {
	nm.getCounter(nid, _readBytesCountTotal, sid).Add(float64(rBytes))
}

// ReportConnWriteBytesCount 链接写入字节数
func ReportConnWriteBytesCount(nid string, sid int64, wBytes int) {
	nm.getCounter(nid, _writeBytesCountTotal, sid).Add(float64(wBytes))
}
