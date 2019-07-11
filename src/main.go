package main

import (
	"fmt"
	"log"
	AppCenter "ms-appcenter-exporter/src/appcenter-client"
	"net/http"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	Author  = "Karin Berg"
	Version = "0.1.0"
)

var (
	argparser       *flags.Parser
	args            []string
	Verbose         bool
	Logger          *DaemonLogger
	AppCenterClient *AppCenter.AppCenterClient

	collectorGeneralList map[string]*CollectorGeneral
	collectorAppList     map[string]*CollectorApp
)

var opts struct {
	// general settings
	Verbose []bool `   long:"verbose" short:"v"                   env:"VERBOSE"                              description:"Verbose mode"`

	// server settings
	ServerBind string `long:"bind"                                env:"SERVER_BIND"                          description:"Server address"                                    default:":8080"`

	// scrape time settings
	ScrapeTime           time.Duration  `long:"scrape.time"                  env:"SCRAPE_TIME"               description:"Default scrape time (time.duration)"                       default:"30m"`
	ScrapeTimeApps       *time.Duration `long:"scrape.time.apps"             env:"SCRAPE_TIME_APPS"          description:"Scrape time for apps metrics (time.duration)"`
	ScrapeTimeUiTestRuns *time.Duration `long:"scrape.time.uitestruns"       env:"SCRAPE_TIME_UITESTRUNS"    description:"Scrape time for ui test run metrics (time.duration)"`
	ScrapeTimeLive       *time.Duration `long:"scrape.time.live"             env:"SCRAPE_TIME_LIVE"          description:"Scrape time for live metrics (time.duration)"              default:"30s"`

	// ignore settings
	AppCenterFilterApps    []string `long:"whitelist.apps"    env:"APPCENTER_FILTER_APPS"    env-delim:" "   description:"Filter apps (slut name)"`
	AppCenterBlacklistApps []string `long:"blacklist.apps"    env:"APPCENTER_BLACKLIST_APPS" env-delim:" "   description:"Filter apps (slut name)"`

	// App Center settings
	AppCenterApiURL       *string `long:"appcenter.api.url"                  env:"APPCENTER_API_URL"         description:"AppCenter API url (empty if hosted by microsoft)"`
	AppCenterAccessToken  string  `long:"appcenter.access-token"             env:"APPCENTER_ACCESS_TOKEN"    description:"AppCenter access token" required:"true"`
	AppCenterOrganisation string  `long:"appcenter.organisation"             env:"APPCENTER_ORGANISATION"    description:"AppCenter organization" required:"true"`

	RequestConcurrencyLimit int64 `long:"request.concurrency"                env:"REQUEST_CONCURRENCY"       description:"Number of concurrent requests against appcenter.ms"  default:"10"`
	RequestRetries          int   `long:"request.retries"                    env:"REQUEST_RETRIES"           description:"Number of retried requests against appcenter.ms"     default:"3"`
}

func main() {
	initArgumentParser()

	// set verbosity
	Verbose = len(opts.Verbose) >= 1

	Logger = NewLogger(log.Lshortfile, Verbose)
	defer Logger.Close()

	Logger.Infof("Init AppCenter exporter v%s (written by %v)", Version, Author)

	Logger.Infof("Init AppCenter connection")
	initAzureConnection()

	Logger.Info("Starting metrics collection")
	Logger.Infof("set scape interval[Default]: %v", scrapeIntervalStatus(&opts.ScrapeTime))
	Logger.Infof("set scape interval[Live]: %v", scrapeIntervalStatus(opts.ScrapeTimeLive))
	Logger.Infof("set scape interval[Apps]: %v", scrapeIntervalStatus(opts.ScrapeTimeApps))
	Logger.Infof("set scape interval[UiTestRuns]: %v", scrapeIntervalStatus(opts.ScrapeTimeUiTestRuns))
	initMetricCollector()

	Logger.Infof("Starting http server on %s", opts.ServerBind)
	startHttpServer()
}

// init argument parser and parse/validate arguments
func initArgumentParser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	// use default scrape time if null
	if opts.ScrapeTimeApps == nil {
		opts.ScrapeTimeApps = &opts.ScrapeTime
	}

	if opts.ScrapeTimeUiTestRuns == nil {
		opts.ScrapeTimeUiTestRuns = &opts.ScrapeTime
	}

	if opts.ScrapeTimeLive == nil {
		opts.ScrapeTimeLive = &opts.ScrapeTime
	}
}

// Init and build AppCenter authorzier
func initAzureConnection() {
	AppCenterClient = AppCenter.NewAppCenterClient()
	if opts.AppCenterApiURL != nil {
		AppCenterClient.HostUrl = opts.AppCenterApiURL
	}

	AppCenterClient.SetOrganization(opts.AppCenterOrganisation)
	AppCenterClient.SetAccessToken(opts.AppCenterAccessToken)
	AppCenterClient.SetConcurrency(opts.RequestConcurrencyLimit)
	AppCenterClient.SetRetries(opts.RequestRetries)
}

func getAppCenterApps() (list AppCenter.AppList) {
	rawList, err := AppCenterClient.ListApps()

	if err != nil {
		panic(err)
	}

	list = rawList

	// whitelist
	if len(opts.AppCenterFilterApps) > 0 {
		rawList = list
		list = AppCenter.AppList{}
		for _, app := range rawList.List {
			if arrayStringContains(opts.AppCenterFilterApps, app.Id) {
				list.List = append(list.List, app)
			}
		}
	}

	// blacklist
	if len(opts.AppCenterBlacklistApps) > 0 {
		// filter ignored appcenter apps
		rawList = list
		list = AppCenter.AppList{}
		for _, app := range rawList.List {
			if !arrayStringContains(opts.AppCenterBlacklistApps, app.Id) {
				list.List = append(list.List, app)
			}
		}
	}

	return
}

func initMetricCollector() {
	var collectorName string
	collectorGeneralList = map[string]*CollectorGeneral{}
	collectorAppList = map[string]*CollectorApp{}

	appList := getAppCenterApps()

	collectorName = "General"
	if opts.ScrapeTimeLive.Seconds() > 0 {
		collectorGeneralList[collectorName] = NewCollectorGeneral(collectorName, &MetricsCollectorGeneral{})
		collectorGeneralList[collectorName].SetAppCenterApps(&appList)
		collectorGeneralList[collectorName].Run(*opts.ScrapeTimeLive)
	} else {
		Logger.Infof("collector[%s]: disabled", collectorName)
	}

	collectorName = "App"
	if opts.ScrapeTimeLive.Seconds() > 0 {
		collectorAppList[collectorName] = NewCollectorApp(collectorName, &MetricsCollectorApp{})
		collectorAppList[collectorName].SetAppCenterApps(&appList)
		collectorAppList[collectorName].Run(*opts.ScrapeTimeLive)
	} else {
		Logger.Infof("collector[%s]: disabled", collectorName)
	}

	collectorName = "UiTestRuns"
	if opts.ScrapeTimeLive.Seconds() > 0 {
		collectorAppList[collectorName] = NewCollectorApp(collectorName, &MetricsCollectorUiTestRuns{})
		collectorAppList[collectorName].SetAppCenterApps(&appList)
		collectorAppList[collectorName].Run(*opts.ScrapeTimeUiTestRuns)
	} else {
		Logger.Infof("collector[%s]: disabled", collectorName)
	}

	collectorName = "LatestUiTestRuns"
	if opts.ScrapeTimeLive.Seconds() > 0 {
		collectorAppList[collectorName] = NewCollectorApp(collectorName, &MetricsCollectorLatestUiTestRun{})
		collectorAppList[collectorName].SetAppCenterApps(&appList)
		collectorAppList[collectorName].Run(*opts.ScrapeTimeUiTestRuns)
	} else {
		Logger.Infof("collector[%s]: disabled", collectorName)
	}

	// background auto update of apps
	if opts.ScrapeTimeApps.Seconds() > 0 {
		go func() {
			// initial sleep
			time.Sleep(*opts.ScrapeTimeApps)

			for {
				Logger.Info("daemon: updating app list")

				appList := getAppCenterApps()

				for _, collector := range collectorGeneralList {
					collector.SetAppCenterApps(&appList)
				}

				for _, collector := range collectorAppList {
					collector.SetAppCenterApps(&appList)
				}

				Logger.Infof("daemon: found %v apps", appList.Count)
				time.Sleep(*opts.ScrapeTimeApps)
			}
		}()
	}
}

// start and handle prometheus handler
func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	Logger.Error(http.ListenAndServe(opts.ServerBind, nil))
}
