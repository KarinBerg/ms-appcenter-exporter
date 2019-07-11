package main

import (
	"context"
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectorApp struct {
	CollectorProcessorApp

	prometheus struct {
		app *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorApp) Setup(collector *CollectorApp) {
	m.CollectorReference = collector

	m.prometheus.app = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "appcenter_app_info",
			Help: "AppCenter apps",
		},
		[]string{
			"appID",
			"appName",
		},
	)

	prometheus.MustRegister(m.prometheus.app)
}

func (m *MetricsCollectorApp) Reset() {
	m.prometheus.app.Reset()
}

func (m *MetricsCollectorApp) Collect(ctx context.Context, callback chan<- func(), app appCenterClient.App) {
	m.collectApp(ctx, callback, app)
}

func (m *MetricsCollectorApp) collectApp(ctx context.Context, callback chan<- func(), app appCenterClient.App) {
	appMetric := MetricCollectorList{}

	appMetric.AddInfo(prometheus.Labels{
		"appID":   app.Id,
		"appName": app.Name,
	})

	callback <- func() {
		appMetric.GaugeSet(m.prometheus.app)
	}
}
