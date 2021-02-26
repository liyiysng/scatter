package metrics

import (
	"runtime"
	"time"
)

const (
	sysSubsystem = "sys"
	goroutines   = "goroutines"
	heapsize     = "heapsize"
	heapobjects  = "heapobjects"
)

// NewSysMetricsSpec 获取系统相关指标
func NewSysMetricsSpec() (*CustomMetricsSpec, error) {
	var spec CustomMetricsSpec

	g := []*Gauge{
		{
			Subsystem: sysSubsystem,
			Name:      goroutines,
			Help:      "the current number of goroutines",
		},
		{
			Subsystem: sysSubsystem,
			Name:      heapsize,
			Help:      "the current heap size",
		},
		{
			Subsystem: sysSubsystem,
			Name:      heapobjects,
			Help:      "the current number of allocated heap objects",
		},
	}

	spec.Gauges = g

	return &spec, nil
}

// ReportSysMetrics reports sys metrics
func ReportSysMetrics(reporters []Reporter, period time.Duration) {
	for {
		for _, r := range reporters {
			num := runtime.NumGoroutine()
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)

			err := r.ReportGauge(goroutines, map[string]string{}, float64(num))
			if err != nil {
				myLog.Error(err)
				return
			}
			r.ReportGauge(heapsize, map[string]string{}, float64(m.Alloc))
			if err != nil {
				myLog.Error(err)
				return
			}
			r.ReportGauge(heapobjects, map[string]string{}, float64(m.HeapObjects))
			if err != nil {
				myLog.Error(err)
				return
			}
		}

		time.Sleep(period)
	}
}
