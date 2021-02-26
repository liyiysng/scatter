// Package metrics 使用 prometheus 检测指标
package metrics

import (
	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("metrics")
)

// Reporter interface
type Reporter interface {
	ReportCount(metric string, tags map[string]string, count float64) error
	ReportSummary(metric string, tags map[string]string, value float64) error
	ReportGauge(metric string, tags map[string]string, value float64) error
	ReportGaugeInc(metric string, tags map[string]string) error
	ReportGaugeDec(metric string, tags map[string]string) error
	Close()
}
