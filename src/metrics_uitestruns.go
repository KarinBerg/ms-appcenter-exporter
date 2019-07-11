package main

import (
	"context"
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectorUiTestRuns struct {
	CollectorProcessorApp

	prometheus struct {
		uiTestRuns *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorUiTestRuns) Setup(collector *CollectorApp) {
	m.CollectorReference = collector

	m.prometheus.uiTestRuns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "appcenter_uitestruns_info",
			Help: "AppCenter UiTest runs",
		},
		[]string{
			"appName",
			"uiTestRunId",
			"uiTestRunDate",
			"uiTestRunSeries",
			"uiTestRunStatus",
			"uiTestRunResultStatus",
			"uiTestRunUrl",
		},
	)

	prometheus.MustRegister(m.prometheus.uiTestRuns)
}

func (m *MetricsCollectorUiTestRuns) Reset() {
	m.prometheus.uiTestRuns.Reset()
}

func (m *MetricsCollectorUiTestRuns) Collect(ctx context.Context, callback chan<- func(), app appCenterClient.App) {
	m.collectUITest(ctx, callback, app)
}

func (m *MetricsCollectorUiTestRuns) collectUITest(ctx context.Context, callback chan<- func(), app appCenterClient.App) {

	uiTestRuns, error := AppCenterClient.ListUiTestRuns(app.Name)
	if error != nil {
		Logger.Errorf("app[%v]call[ListUiTestRuns]: %v", app.Name, error)
		return
	}

	uiTestRunMetric := MetricCollectorList{}

	for _, uiTestRun := range uiTestRuns.List {

		uiTestRunMetric.AddInfo(prometheus.Labels{
			"appName":               uiTestRuns.AppName,
			"uiTestRunId":           uiTestRun.ID,
			"uiTestRunDate":         uiTestRun.Date.Format(time.RFC3339),
			"uiTestRunSeries":       uiTestRun.TestSeries,
			"uiTestRunStatus":       uiTestRun.RunStatus,
			"uiTestRunResultStatus": uiTestRun.ResultStatus,
			"uiTestRunUrl":          "https://appcenter.ms/orgs/" + url.QueryEscape(AppCenterClient.GetOrganization()) + "/apps/" + url.QueryEscape(uiTestRuns.AppName) + "/test/runs/" + uiTestRun.ID,
		})
	}

	callback <- func() {
		uiTestRunMetric.GaugeSet(m.prometheus.uiTestRuns)
	}
}
