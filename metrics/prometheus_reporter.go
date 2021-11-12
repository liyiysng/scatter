package metrics

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/liyiysng/scatter/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/viper"
)

// PrometheusReporter Prometheus 监控指标
type PrometheusReporter struct {
	countReportersMap   map[string]*prometheus.CounterVec
	summaryReportersMap map[string]*prometheus.SummaryVec
	gaugeReportersMap   map[string]*prometheus.GaugeVec
	additionalLabels    map[string]string

	pusher     *push.Pusher
	closeEvent *util.Event
	wg         sync.WaitGroup
}

// Close 关闭
func (p *PrometheusReporter) Close() {
	p.closeEvent.Fire()
	p.wg.Wait()
}

func (p *PrometheusReporter) registerMetrics(
	ns string,
	constLabels map[string]string,
	additionalLabelsKeys []string,
	spec *CustomMetricsSpec,
) {
	for _, summary := range spec.Summaries {
		p.summaryReportersMap[summary.Name] = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   ns,
				Subsystem:   summary.Subsystem,
				Name:        summary.Name,
				Help:        summary.Help,
				Objectives:  summary.Objectives,
				ConstLabels: constLabels,
			},
			append(additionalLabelsKeys, summary.Labels...),
		)
	}

	for _, gauge := range spec.Gauges {
		p.gaugeReportersMap[gauge.Name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   ns,
				Subsystem:   gauge.Subsystem,
				Name:        gauge.Name,
				Help:        gauge.Help,
				ConstLabels: constLabels,
			},
			append(additionalLabelsKeys, gauge.Labels...),
		)
	}

	for _, counter := range spec.Counters {
		p.countReportersMap[counter.Name] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace:   ns,
				Subsystem:   counter.Subsystem,
				Name:        counter.Name,
				Help:        counter.Help,
				ConstLabels: constLabels,
			},
			append(additionalLabelsKeys, counter.Labels...),
		)
	}
}

func (p *PrometheusReporter) registerAllMetrics(
	config *viper.Viper,
) error {

	additionalLabels := config.GetStringMapString("scatter.metrics.additionalTags")
	p.additionalLabels = additionalLabels
	ns := config.GetString("scatter.metrics.prometheus.namespace")

	additionalLabelsKeys := make([]string, 0, len(additionalLabels))
	for key := range additionalLabels {
		additionalLabelsKeys = append(additionalLabelsKeys, key)
	}

	constLabels := config.GetStringMapString("scatter.metrics.constTags")

	// 系统指标
	spec, err := NewSysMetricsSpec()
	if err != nil {
		return err
	}
	p.registerMetrics(ns, constLabels, additionalLabelsKeys, spec)

	// 节点指标
	spec, err = NewNodeMetricsSpec()
	if err != nil {
		return err
	}
	p.registerMetrics(ns, constLabels, additionalLabelsKeys, spec)

	// 自定义指标
	spec, err = NewCustomMetricsSpec(config)
	if err != nil {
		return err
	}
	p.registerMetrics(ns, constLabels, additionalLabelsKeys, spec)

	toRegister := make([]prometheus.Collector, 0)
	for _, c := range p.countReportersMap {
		toRegister = append(toRegister, c)
	}

	for _, c := range p.gaugeReportersMap {
		toRegister = append(toRegister, c)
	}

	for _, c := range p.summaryReportersMap {
		toRegister = append(toRegister, c)
	}

	if p.pusher != nil {
		registry := prometheus.NewRegistry()
		registry.MustRegister(toRegister...)
		p.pusher.Gatherer(registry)
	} else {
		prometheus.MustRegister(toRegister...)
	}

	return nil
}

// NewPrometheusReporter 创建
func NewPrometheusReporter(
	config *viper.Viper,
) (Reporter, error) {

	prometheusReporter := &PrometheusReporter{
		closeEvent:          util.NewEvent(),
		countReportersMap:   make(map[string]*prometheus.CounterVec),
		summaryReportersMap: make(map[string]*prometheus.SummaryVec),
		gaugeReportersMap:   make(map[string]*prometheus.GaugeVec),
	}

	if config.GetString("scatter.metrics.prometheus.collect_type") == "listen" {
		port := config.GetInt("scatter.metrics.prometheus.port")
		prometheusReporter.registerAllMetrics(config)
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			lErr := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if lErr != nil {
				myLog.Fatal(lErr)
			}
		}()
	} else {
		jobName := config.GetString("scatter.metrics.prometheus.pusher.job_name")
		prometheusReporter.pusher = push.New(config.GetString("scatter.metrics.prometheus.pusher.addr"), jobName)
		prometheusReporter.registerAllMetrics(config)
		err := prometheusReporter.pusher.Push()
		if err != nil {
			return nil, err
		}
		go func() {

			waitTicker := time.NewTicker(config.GetDuration("scatter.metrics.prometheus.pusher.inteval"))
			defer waitTicker.Stop()

			prometheusReporter.wg.Add(1)

			defer func() {
				delErr := prometheusReporter.pusher.Delete()
				if delErr != nil {
					myLog.Errorf("prometheusReporter.pusher delete error %v", delErr)
				}
				prometheusReporter.wg.Done()
			}()

			for {

				select {
				case <-waitTicker.C:
					{
						addErr := prometheusReporter.pusher.Add()
						if addErr != nil {
							myLog.Errorf("prometheusReporter.pusher error %v", err)
						}
					}
				case <-prometheusReporter.closeEvent.Done():
					{
						return
					}
				}
			}

		}()
	}
	return prometheusReporter, nil
}

// ReportSummary reports a summary metric
func (p *PrometheusReporter) ReportSummary(metric string, labels map[string]string, value float64) error {

	if p.closeEvent.HasFired() {
		return ErrReporterClosed
	}

	sum := p.summaryReportersMap[metric]
	if sum != nil {
		labels = p.ensureLabels(labels)
		sum.With(labels).Observe(value)
		return nil
	}
	return ErrMetricNotKnown
}

// ReportCount reports a summary metric
func (p *PrometheusReporter) ReportCount(metric string, labels map[string]string, count float64) error {

	if p.closeEvent.HasFired() {
		return ErrReporterClosed
	}

	cnt := p.countReportersMap[metric]
	if cnt != nil {
		labels = p.ensureLabels(labels)
		cnt.With(labels).Add(count)
		return nil
	}
	return ErrMetricNotKnown
}

// ReportGauge reports a gauge metric
func (p *PrometheusReporter) ReportGauge(metric string, labels map[string]string, value float64) error {

	if p.closeEvent.HasFired() {
		return ErrReporterClosed
	}

	g := p.gaugeReportersMap[metric]
	if g != nil {
		labels = p.ensureLabels(labels)
		g.With(labels).Set(value)
		return nil
	}
	return ErrMetricNotKnown
}

// ReportGaugeInc reports a gauge metric
func (p *PrometheusReporter) ReportGaugeInc(metric string, labels map[string]string) error {

	if p.closeEvent.HasFired() {
		return ErrReporterClosed
	}

	g := p.gaugeReportersMap[metric]
	if g != nil {
		labels = p.ensureLabels(labels)
		g.With(labels).Inc()
		return nil
	}
	return ErrMetricNotKnown
}

// ReportGaugeDec reports a gauge metric
func (p *PrometheusReporter) ReportGaugeDec(metric string, labels map[string]string) error {

	if p.closeEvent.HasFired() {
		return ErrReporterClosed
	}

	g := p.gaugeReportersMap[metric]
	if g != nil {
		labels = p.ensureLabels(labels)
		g.With(labels).Dec()
		return nil
	}
	return ErrMetricNotKnown
}

// ensureLabels checks if labels contains the additionalLabels values,
// otherwise adds them with the default values
func (p *PrometheusReporter) ensureLabels(labels map[string]string) map[string]string {
	for key, defaultVal := range p.additionalLabels {
		if _, ok := labels[key]; !ok {
			labels[key] = defaultVal
		}
	}

	return labels
}
