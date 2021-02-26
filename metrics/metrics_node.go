package metrics

import (
	"strconv"
	"time"
)

const (
	nodeSubSystem = "node"
	// 当前链接数
	connCount = "conn_count"
	// 收发字节
	readBytesTotal  = "read_bytes_total"
	writeBytesTotal = "write_bytes_total"
	// 消息处理延时
	msgProcDelay = "msg_proc_delay"
)

// NewNodeMetricsSpec 获取节点相关指标
func NewNodeMetricsSpec() (*CustomMetricsSpec, error) {
	var spec CustomMetricsSpec

	spec.Gauges = []*Gauge{
		{
			Subsystem: nodeSubSystem,
			Name:      connCount,
			Help:      "the current connection count of node",
			Labels:    []string{"nid", "node_name"},
		},
	}

	spec.Counters = []*Counter{
		{
			Subsystem: nodeSubSystem,
			Name:      readBytesTotal,
			Help:      "the current read bytes total of node",
			Labels:    []string{"nid", "node_name"},
		},
		{
			Subsystem: nodeSubSystem,
			Name:      writeBytesTotal,
			Help:      "the current written bytes total of node",
			Labels:    []string{"nid", "node_name"},
		},
	}

	spec.Summaries = []*Summary{
		{
			Subsystem: nodeSubSystem,
			Name:      msgProcDelay,
			Help:      "the message process delay(ms) of node",
			Labels:    []string{"nid", "node_name", "srv", "method"},
		},
	}

	return &spec, nil
}

// ReportNodeConnectionInc 当前链接数加一
func ReportNodeConnectionInc(reporters []Reporter, nid int64, nname string) {
	for _, r := range reporters {
		err := r.ReportGaugeInc(connCount, map[string]string{"nid": strconv.FormatInt(nid, 10), "node_name": nname})
		if err != nil {
			myLog.Error(err)
			return
		}
	}
}

// ReportNodeConnectionDec 当前链接数减一
func ReportNodeConnectionDec(reporters []Reporter, nid int64, nname string) {
	for _, r := range reporters {
		err := r.ReportGaugeDec(connCount, map[string]string{"nid": strconv.FormatInt(nid, 10), "node_name": nname})
		if err != nil {
			myLog.Error(err)
			return
		}
	}
}

// ReportNodeReadBytes 读字节数
func ReportNodeReadBytes(reporters []Reporter, nid int64, nname string, count int) {
	for _, r := range reporters {
		err := r.ReportCount(readBytesTotal, map[string]string{"nid": strconv.FormatInt(nid, 10), "node_name": nname}, float64(count))
		if err != nil {
			myLog.Error(err)
			return
		}
	}
}

// ReportNodeWriteBytes 写字节数
func ReportNodeWriteBytes(reporters []Reporter, nid int64, nname string, count int) {
	for _, r := range reporters {
		err := r.ReportCount(writeBytesTotal, map[string]string{"nid": strconv.FormatInt(nid, 10), "node_name": nname}, float64(count))
		if err != nil {
			myLog.Error(err)
			return
		}
	}
}

// ReportMsgProcDelay 消息处理延时
func ReportMsgProcDelay(reporters []Reporter, nid int64, nname string, srv, method string, delay time.Duration) {
	for _, r := range reporters {
		err := r.ReportSummary(msgProcDelay, map[string]string{"nid": strconv.FormatInt(nid, 10), "node_name": nname, "srv": srv, "method": method}, float64(delay.Milliseconds()))
		if err != nil {
			myLog.Error(err)
			return
		}
	}
}
