package AppCenterClient

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/go-resty/resty"
)

type AppCenterClient struct {
	organization *string
	accessToken  *string

	HostUrl *string

	restClient *resty.Client

	semaphore   chan bool
	concurrency int64

	RequestCount   uint64
	RequestRetries int
}

func NewAppCenterClient() *AppCenterClient {
	c := AppCenterClient{}
	c.Init()

	return &c
}

func (c *AppCenterClient) Init() {
	c.RequestCount = 0
	c.SetRetries(3)
	c.SetConcurrency(10)
}

func (c *AppCenterClient) SetConcurrency(v int64) {
	c.concurrency = v
	c.semaphore = make(chan bool, c.concurrency)
}

func (c *AppCenterClient) SetRetries(v int) {
	c.RequestRetries = v

	if c.restClient != nil {
		c.restClient.SetRetryCount(c.RequestRetries)
	}
}

func (c *AppCenterClient) SetOrganization(url string) {
	c.organization = &url
}

func (c *AppCenterClient) GetOrganization() string {
	return *c.organization
}

func (c *AppCenterClient) SetAccessToken(token string) {
	c.accessToken = &token
}

func (c *AppCenterClient) rest() *resty.Client {
	if c.restClient == nil {
		c.restClient = resty.New()
		if c.HostUrl != nil {
			c.restClient.SetHostURL(*c.HostUrl)
		} else {
			c.restClient.SetHostURL("https://api.appcenter.ms/v0.1/")
		}
		c.restClient.SetHeader("Accept", "application/json")
		c.restClient.SetHeader("X-API-Token", *c.accessToken)
		//c.restClient.SetBasicAuth("", *c.accessToken)
		c.restClient.SetRetryCount(c.RequestRetries)
		c.restClient.OnBeforeRequest(c.restOnBeforeRequest)
		c.restClient.OnAfterResponse(c.restOnAfterResponse)
	}

	return c.restClient
}

func (c *AppCenterClient) concurrencyLock() {
	c.semaphore <- true
}

func (c *AppCenterClient) concurrencyUnlock() {
	<-c.semaphore
}

func (c *AppCenterClient) restOnBeforeRequest(client *resty.Client, request *resty.Request) (err error) {
	atomic.AddUint64(&c.RequestCount, 1)
	return
}

func (c *AppCenterClient) restOnAfterResponse(client *resty.Client, response *resty.Response) (err error) {
	return
}

func (c *AppCenterClient) GetRequestCount() float64 {
	requestCount := atomic.LoadUint64(&c.RequestCount)
	return float64(requestCount)
}

func (c *AppCenterClient) GetCurrentConcurrency() float64 {
	return float64(len(c.semaphore))
}

func (c *AppCenterClient) checkResponse(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response != nil {
		// check status code
		statusCode := response.StatusCode()
		if statusCode != 200 {
			return errors.New(fmt.Sprintf("Response status code is %v (expected 200)", statusCode))
		}
	} else {
		return errors.New("Response is nil")
	}

	return nil
}
