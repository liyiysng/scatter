package metrics

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/liyiysng/scatter/config"
	"github.com/liyiysng/scatter/constants"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	prometheusReporter *PrometheusReporter
	once               sync.Once
)

// PrometheusReporter Prometheus 监控指标
type PrometheusReporter struct {
	serverType          string
	game                string
	countReportersMap   map[string]*prometheus.CounterVec
	summaryReportersMap map[string]*prometheus.SummaryVec
	gaugeReportersMap   map[string]*prometheus.GaugeVec
	additionalLabels    map[string]string
}

func (p *PrometheusReporter) registerCustomMetrics(
	constLabels map[string]string,
	additionalLabelsKeys []string,
	spec *CustomMetricsSpec,
) {
	for _, summary := range spec.Summaries {
		p.summaryReportersMap[summary.Name] = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:   "scatter",
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
				Namespace:   "scatter",
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
				Namespace:   "scatter",
				Subsystem:   counter.Subsystem,
				Name:        counter.Name,
				Help:        counter.Help,
				ConstLabels: constLabels,
			},
			append(additionalLabelsKeys, counter.Labels...),
		)
	}
}

func (p *PrometheusReporter) registerMetrics(
	constLabels, additionalLabels map[string]string,
	spec *CustomMetricsSpec,
) {

	constLabels["game"] = p.game
	constLabels["serverType"] = p.serverType

	p.additionalLabels = additionalLabels
	additionalLabelsKeys := make([]string, 0, len(additionalLabels))
	for key := range additionalLabels {
		additionalLabelsKeys = append(additionalLabelsKeys, key)
	}

	p.registerCustomMetrics(constLabels, additionalLabelsKeys, spec)

	// HandlerResponseTimeMs summary
	p.summaryReportersMap[ResponseTime] = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:   "scatter",
			Subsystem:   "handler",
			Name:        ResponseTime,
			Help:        "the time to process a msg in nanoseconds",
			Objectives:  map[float64]float64{0.7: 0.02, 0.95: 0.005, 0.99: 0.001},
			ConstLabels: constLabels,
		},
		append([]string{"route", "status", "type", "code"}, additionalLabelsKeys...),
	)

	// ProcessDelay summary
	p.summaryReportersMap[ProcessDelay] = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:   "scatter",
			Subsystem:   "handler",
			Name:        ProcessDelay,
			Help:        "the delay to start processing a msg in nanoseconds",
			Objectives:  map[float64]float64{0.7: 0.02, 0.95: 0.005, 0.99: 0.001},
			ConstLabels: constLabels,
		},
		append([]string{"route", "type"}, additionalLabelsKeys...),
	)

	// ConnectedClients gauge
	p.gaugeReportersMap[ConnectedClients] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "acceptor",
			Name:        ConnectedClients,
			Help:        "the number of clients connected right now",
			ConstLabels: constLabels,
		},
		additionalLabelsKeys,
	)

	p.gaugeReportersMap[Goroutines] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "sys",
			Name:        Goroutines,
			Help:        "the current number of goroutines",
			ConstLabels: constLabels,
		},
		additionalLabelsKeys,
	)

	p.gaugeReportersMap[HeapSize] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "sys",
			Name:        HeapSize,
			Help:        "the current heap size",
			ConstLabels: constLabels,
		},
		additionalLabelsKeys,
	)

	p.gaugeReportersMap[HeapObjects] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "scatter",
			Subsystem:   "sys",
			Name:        HeapObjects,
			Help:        "the current number of allocated heap objects",
			ConstLabels: constLabels,
		},
		additionalLabelsKeys,
	)

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

	prometheus.MustRegister(toRegister...)
}

// GetPrometheusReporter 获取指标监控单例
func GetPrometheusReporter(
	serverType string,
	config *config.Config,
	constLabels map[string]string,
) (*PrometheusReporter, error) {
	port := config.GetInt("scatter.metrics.prometheus.port")
	game := config.GetString("scatter.game")
	additionalLabels := config.GetStringMapString("scatter.metrics.additionalTags")

	spec, err := NewCustomMetricsSpec(config)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		prometheusReporter = &PrometheusReporter{
			serverType:          serverType,
			game:                game,
			countReportersMap:   make(map[string]*prometheus.CounterVec),
			summaryReportersMap: make(map[string]*prometheus.SummaryVec),
			gaugeReportersMap:   make(map[string]*prometheus.GaugeVec),
		}
		prometheusReporter.registerMetrics(constLabels, additionalLabels, spec)
		http.Handle("/metrics", promhttp.Handler())
		go (func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
		})()
	})

	return prometheusReporter, nil
}

// ReportSummary reports a summary metric
func (p *PrometheusReporter) ReportSummary(metric string, labels map[string]string, value float64) error {
	sum := p.summaryReportersMap[metric]
	if sum != nil {
		labels = p.ensureLabels(labels)
		sum.With(labels).Observe(value)
		return nil
	}
	return constants.ErrMetricNotKnown
}

// ReportCount reports a summary metric
func (p *PrometheusReporter) ReportCount(metric string, labels map[string]string, count float64) error {
	cnt := p.countReportersMap[metric]
	if cnt != nil {
		labels = p.ensureLabels(labels)
		cnt.With(labels).Add(count)
		return nil
	}
	return constants.ErrMetricNotKnown
}

// ReportGauge reports a gauge metric
func (p *PrometheusReporter) ReportGauge(metric string, labels map[string]string, value float64) error {
	g := p.gaugeReportersMap[metric]
	if g != nil {
		labels = p.ensureLabels(labels)
		g.With(labels).Set(value)
		return nil
	}
	return constants.ErrMetricNotKnown
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
