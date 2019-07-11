package AppCenterClient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"
)

type UiTestRunList struct {
	Count   int `json:"count"`
	AppName string
	List    []UiTestRun `json:"value"`
}

type UiTestRun struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"date"`
	Platform       string    `json:"platform"`
	Stats          Stats     `json:"stats"`
	RunStatus      string    `json:"runStatus"`
	ResultStatus   string    `json:"resultStatus"`
	State          string    `json:"state"`
	Status         string    `json:"status"`
	Description    string    `json:"description"`
	AppVersion     string    `json:"appVersion"`
	TestSeries     string    `json:"testSeries"`
	TestTypeStatus string    `json:"testType"`
}

type Stats struct {
	Devices            int     `json:"devices"`
	DevicesFinished    int     `json:"devicesFinished"`
	DevicesFailed      int     `json:"devicesFailed"`
	Total              int     `json:"total"`
	Passed             int     `json:"passed"`
	Failed             int     `json:"failed"`
	PeakMemory         float32 `json:"peakMemory"`
	TotalDeviceMinutes int     `json:"totalDeviceMinutes"`
}

func (c *AppCenterClient) ListUiTestRuns(appname string) (list UiTestRunList, error error) {
	defer c.concurrencyUnlock()
	c.concurrencyLock()

	response, err := c.rest().R().Get(fmt.Sprintf("apps/%v/%v/test_runs", url.QueryEscape(*c.organization), url.QueryEscape(appname)))
	if err := c.checkResponse(response, err); err != nil {
		error = err
		return
	}

	var slist []UiTestRun
	err = json.Unmarshal(response.Body(), &slist)
	if err != nil {
		error = err
		return
	}

	list.Count = len(slist)
	list.List = slist
	list.AppName = appname
	return
}

func (c *AppCenterClient) ListLastestUiTestRun(appname string) (test *UiTestRun, error error) {
	defer c.concurrencyUnlock()
	c.concurrencyLock()
	response, err := c.rest().R().Get(fmt.Sprintf("apps/%v/%v/test_runs", url.QueryEscape(*c.organization), url.QueryEscape(appname)))
	if err := c.checkResponse(response, err); err != nil {
		error = err
		return
	}

	var slist []UiTestRun
	err = json.Unmarshal(response.Body(), &slist)
	if err != nil {
		error = err
		return
	}

	if len(slist) <= 0 {
		test = nil
		return
	}
	sort.Slice(slist, func(i, j int) bool { return slist[i].Date.After(slist[j].Date) })
	test = &slist[0]
	return
}
