// Package metrics 使用 prometheus 检测指标
package metrics

import (
	"errors"

	"github.com/liyiysng/scatter/logger"
)

var (
	myLog = logger.Component("metrics")
)

var (
	// ErrMetricNotKnown 未知的指标
	ErrMetricNotKnown = errors.New("the provided metric does not exist")
	// ErrReporterClosed 指标回报关闭
	ErrReporterClosed = errors.New("metric reporter closed")
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
