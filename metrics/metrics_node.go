package metrics

import "strconv"

const (
	nodeSubSystem = "node"
	// 当前链接数
	connCount = "conn_count"
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
