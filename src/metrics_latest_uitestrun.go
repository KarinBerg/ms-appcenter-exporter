package main

import (
	"context"
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollectorLatestUiTestRun struct {
	CollectorProcessorApp

	prometheus struct {
		latestUiTestRun *prometheus.GaugeVec
	}
}

func (m *MetricsCollectorLatestUiTestRun) Setup(collector *CollectorApp) {
	m.CollectorReference = collector

	m.prometheus.latestUiTestRun = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "appcenter_latest_uitestrun_info",
			Help: "AppCenter Latest UiTest Run",
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

	prometheus.MustRegister(m.prometheus.latestUiTestRun)
}

func (m *MetricsCollectorLatestUiTestRun) Reset() {
	m.prometheus.latestUiTestRun.Reset()
}

func (m *MetricsCollectorLatestUiTestRun) Collect(ctx context.Context, callback chan<- func(), app appCenterClient.App) {
	m.collectUITest(ctx, callback, app)
}

func (m *MetricsCollectorLatestUiTestRun) collectUITest(ctx context.Context, callback chan<- func(), app appCenterClient.App) {

	uiTestRun, error := AppCenterClient.ListLastestUiTestRun(app.Name)
	if error != nil {
		Logger.Errorf("app[%v]call[ListLastestUiTestRun]: %v", app.Name, error)
		return
	}

	if uiTestRun == nil {
		return
	}

	uiTestRunMetric := MetricCollectorList{}

	uiTestRunMetric.AddInfo(prometheus.Labels{
		"appName":               app.Name,
		"uiTestRunId":           uiTestRun.ID,
		"uiTestRunDate":         uiTestRun.Date.Format(time.RFC3339),
		"uiTestRunSeries":       uiTestRun.TestSeries,
		"uiTestRunStatus":       uiTestRun.RunStatus,
		"uiTestRunResultStatus": uiTestRun.ResultStatus,
		"uiTestRunUrl":          "https://appcenter.ms/orgs/" + url.QueryEscape(AppCenterClient.GetOrganization()) + "/apps/" + url.QueryEscape(app.Name) + "/test/runs/" + uiTestRun.ID,
	})

	callback <- func() {
		uiTestRunMetric.GaugeSet(m.prometheus.latestUiTestRun)
	}
}
