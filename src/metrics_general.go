package main

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectorGeneral struct {
	CollectorProcessorGeneral

	prometheus struct {
		stats *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorGeneral) Setup(collector *CollectorGeneral) {
	m.CollectorReference = collector

	m.prometheus.stats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "appcenter_stats",
			Help: "App Center statistics",
		},
		[]string{
			"name",
			"type",
		},
	)

	prometheus.MustRegister(m.prometheus.stats)
}

func (m *MetricsCollectorGeneral) Reset() {
	m.prometheus.stats.Reset()
}

func (m *MetricsCollectorGeneral) Collect(ctx context.Context, callback chan<- func()) {
	m.collectAppCenterClientStats(ctx, callback)
	m.collectCollectorStats(ctx, callback)
}

func (m *MetricsCollectorGeneral) collectAppCenterClientStats(ctx context.Context, callback chan<- func()) {
	statsMetrics := MetricCollectorList{}

	statsMetrics.Add(prometheus.Labels{
		"name": "appcenter.ms",
		"type": "requests",
	}, AppCenterClient.GetRequestCount())

	statsMetrics.Add(prometheus.Labels{
		"name": "appcenter.ms",
		"type": "concurrency",
	}, AppCenterClient.GetCurrentConcurrency())

	callback <- func() {
		statsMetrics.GaugeSet(m.prometheus.stats)
	}
}

func (m *MetricsCollectorGeneral) collectCollectorStats(ctx context.Context, callback chan<- func()) {
	statsMetrics := MetricCollectorList{}

	for _, collector := range collectorGeneralList {
		if collector.LastScrapeDuration != nil {
			statsMetrics.AddDuration(prometheus.Labels{
				"name": collector.Name,
				"type": "collectorDuration",
			}, *collector.LastScrapeDuration)
		}
	}

	for _, collector := range collectorAppList {
		if collector.LastScrapeDuration != nil {
			statsMetrics.AddDuration(prometheus.Labels{
				"name": collector.Name,
				"type": "collectorDuration",
			}, *collector.LastScrapeDuration)
		}
	}

	callback <- func() {
		statsMetrics.GaugeSet(m.prometheus.stats)
	}
}
