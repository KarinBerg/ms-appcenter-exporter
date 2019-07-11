package main

import (
	"context"
	appCenterClient "ms-appcenter-exporter/src/appcenter-client"
)

type CollectorProcessorAppInterface interface {
	Setup(collector *CollectorApp)
	Reset()
	Collect(ctx context.Context, callback chan<- func(), app appCenterClient.App)
}

type CollectorProcessorApp struct {
	CollectorProcessorAppInterface
	CollectorReference *CollectorApp
}

func NewCollectorApp(name string, processor CollectorProcessorAppInterface) *CollectorApp {
	collector := CollectorApp{
		CollectorBase: CollectorBase{
			Name: name,
		},
		Processor: processor,
	}

	return &collector
}
