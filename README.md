Microsoft AppCenter Exporter
============================

Microsoft AppCenter Api: https://openapi.appcenter.ms

[![license](https://img.shields.io/github/license/KarinBerg/ms-appcenter-exporter.svg)](https://github.com/KarinBerg/ms-appcenter-exporter/blob/master/LICENSE)
[![Docker](https://img.shields.io/badge/docker-dockerberg%2Fms--appcenter--exporter-blue.svg?longCache=true&style=flat&logo=docker)](https://hub.docker.com/r/dockerberg/ms-appcenter-exporter/)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/dockerberg/ms-appcenter-exporter.svg)](https://hub.docker.com/r/dockerberg/ms-appcenter-exporter/builds)

Prometheus exporter for Microsoft AppCenter for apps info and ui test runs.

Configuration
-------------

Normally no configuration is needed but can be customized using environment variables.

| Environment variable                  | DefaultValue                        | Description                                                              |
|---------------------------------------|-------------------------------------|--------------------------------------------------------------------------|
| `SCRAPE_TIME`                         | `30m`                               | Interval (time.Duration) between API calls                               |
| `SCRAPE_TIME_APPS`                    | not set, default see `SCRAPE_TIME`  | Interval for app metrics (list of apps for all scrapers)         |
| `SCRAPE_TIME_UITESTRUNS`              | not set, default see `SCRAPE_TIME`  | Interval for ui test runs metrics                                          |
| `SCRAPE_TIME_LIVE`                    | `30s`                               | Time (time.Duration) between API calls                                   |
| `SERVER_BIND`                         | `:8080`                             | IP/Port binding                                                          |
| `APPCENTER_API_URL`                   | none                                | MS AppCenter API url (only if on-prem)                                       |
| `APPCENTER_ORGANISATION`              | none                                | AppCenter organisation            |
| `APPCENTER_ACCESS_TOKEN`              | none                                | AppCenter access token                                                |
| `APPCENTER_FILTER_APPS`               | none                                | Whitelist project uuids                                                  |
| `APPCENTER_BLACKLIST_APPS`            | none                                | Blacklist project uuids                                                  |
| `REQUEST_CONCURRENCY`                 | `10`                                | API request concurrency (number of calls at the same time)              |
| `REQUEST_RETRIES`                     | `3`                                 | API request retries in case of failure                                 |


Metrics
-------

| Metric                                          | Scraper       | Description                                                                          |
|-------------------------------------------------|---------------|--------------------------------------------------------------------------------------|
| `appcenter_stats`                               | live          | General scraper stats                                                                |
| `appcenter_app_info`                            | live          | Organization app information                                                        |
| `appcenter_uitestruns_info`                     | live          | Count of ui test runs (by status)                                                          |
| `appcenter_latest_uitestrun_info`               | live          | Latest ui test run status informations                                                     |

Usage
-----

```
Usage:
  ms-appcenter-exporter [OPTIONS]

Application Options:
  -v, --verbose                      Verbose mode [$VERBOSE]
      --bind=                        Server address (default: :8080) [$SERVER_BIND]
      --scrape.time=                 Default scrape time (time.duration) (default: 30m) [$SCRAPE_TIME]
      --scrape.time.apps=            Scrape time for apps metrics (time.duration) [$SCRAPE_TIME_PROJECTS]
      --scrape.time.uitestruns=      Scrape time for uitestruns metrics (time.duration) [$SCRAPE_TIME_REPOSITORY]
      --scrape.time.live=            Scrape time for live metrics (time.duration) (default: 30s) [$SCRAPE_TIME_LIVE]
      --whitelist.apps=              Filter apps (UUIDs) [$APPCENTER_FILTER_PROJECT]
      --blacklist.apps=              Filter apps (UUIDs) [$APPCENTER_BLACKLIST_PROJECT]
      --appcenter.access-token=      AppCenter access token [$APPCENTER_ACCESS_TOKEN]
      --appcenter.organisation=      AppCenter organization [$APPCENTER_ORGANISATION]
      --request.concurrency=         Number of concurrent requests against api.appcenter.ms (default:10) [$REQUEST_CONCURRENCY]
      --request.retries=             Number of retried requests against api.appcenter.ms (default: 3) [$REQUEST_RETRIES]

Help Options:
  -h, --help                         Show this help message
```
