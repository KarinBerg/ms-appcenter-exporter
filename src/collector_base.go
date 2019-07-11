package main

import (
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"
	"sync"
	"time"
)

type CollectorBase struct {
	Name       string
	scrapeTime *time.Duration

	appCenterApps      *appCenterClient.AppList
	appCenterAppsMutex sync.Mutex

	LastScrapeDuration  *time.Duration
	collectionStartTime time.Time
}

func (c *CollectorBase) Init() {
}

func (c *CollectorBase) SetScrapeTime(scrapeTime time.Duration) {
	c.scrapeTime = &scrapeTime
}

func (c *CollectorBase) GetScrapeTime() *time.Duration {
	return c.scrapeTime
}

func (c *CollectorBase) SetAppCenterApps(apps *appCenterClient.AppList) {
	c.appCenterAppsMutex.Lock()
	c.appCenterApps = apps
	c.appCenterAppsMutex.Unlock()
}

func (c *CollectorBase) GetAppCenterApps() (apps *appCenterClient.AppList) {
	c.appCenterAppsMutex.Lock()
	apps = c.appCenterApps
	c.appCenterAppsMutex.Unlock()
	return
}

func (c *CollectorBase) collectionStart() {
	c.collectionStartTime = time.Now()

	Logger.Infof("collector[%s]: starting metrics collection", c.Name)
}

func (c *CollectorBase) collectionFinish() {
	duration := time.Now().Sub(c.collectionStartTime)
	c.LastScrapeDuration = &duration

	Logger.Infof("collector[%s]: finished metrics collection (duration: %v)", c.Name, c.LastScrapeDuration)
}

func (c *CollectorBase) sleepUntilNextCollection() {
	Logger.Verbosef("collector[%s]: sleeping %v", c.Name, c.GetScrapeTime().String())
	time.Sleep(*c.GetScrapeTime())
}
