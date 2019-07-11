package main

import (
	"context"
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"
	"sync"
	"time"
)

type CollectorApp struct {
	CollectorBase
	Processor CollectorProcessorAppInterface
}

func (c *CollectorApp) Run(scrapeTime time.Duration) {
	c.SetScrapeTime(scrapeTime)

	c.Processor.Setup(c)
	go func() {
		for {
			go func() {
				c.Collect()
			}()
			c.sleepUntilNextCollection()
		}
	}()
}

func (c *CollectorApp) Collect() {
	var wg sync.WaitGroup
	var wgCallback sync.WaitGroup

	if c.GetAppCenterApps() == nil {
		Logger.Infof("collector[%s]: no apps found, skipping", c.Name)
		return
	}

	ctx := context.Background()

	callbackChannel := make(chan func())

	c.collectionStart()

	for _, app := range c.GetAppCenterApps().List {
		wg.Add(1)
		go func(ctx context.Context, callback chan<- func(), app appCenterClient.App) {
			defer wg.Done()
			c.Processor.Collect(ctx, callbackChannel, app)
		}(ctx, callbackChannel, app)
	}

	// collect metrics (callbacks) and proceses them
	wgCallback.Add(1)
	go func() {
		defer wgCallback.Done()
		var callbackList []func()
		for callback := range callbackChannel {
			callbackList = append(callbackList, callback)
		}

		// reset metric values
		c.Processor.Reset()

		// process callbacks (set metrics)
		for _, callback := range callbackList {
			callback()
		}
	}()

	// wait for all funcs
	wg.Wait()
	close(callbackChannel)
	wgCallback.Wait()

	c.collectionFinish()
}
