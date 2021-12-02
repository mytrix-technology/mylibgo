package client

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	defaultWaitTime = 3 * time.Second
	defaultMaxWaitTime = 10 * time.Second
)

type Option func(options *clientOptions)

type clientOptions struct {
	RetryCount int
	RetryWaitTime time.Duration
	KeepAlive bool
}

type HttpClient struct {
	http.Client
	method string
	url string
	requestBody		[]byte
	contentType		HttpCallContentType
	headers			http.Header
	//options clientOptions
	form url.Values
}

func NewWithTimeout(dialTimeout time.Duration, requestTimeout time.Duration) *HttpClient {
	netTransport := &http.Transport{
		DialContext: (&net.Dialer{Timeout: dialTimeout}).DialContext,
		TLSHandshakeTimeout: dialTimeout,
	}

	client := new(HttpClient)
	client.method = "GET"
	client.contentType = HTTP_CALL_CONTENT_JSON

	client.Timeout = requestTimeout
	client.Transport = netTransport

	client.headers = http.Header{}
	client.headers.Set("Content-Type", string(client.contentType))

	return client
}

func (c *HttpClient) Method(httpMethod string) *HttpClient {
	client := *c
	client.method = strings.ToUpper(httpMethod)
	return &client
}

func (c *HttpClient) URL(url string) *HttpClient {
	client := *c
	client.url = url
	return &client
}

//func (c *HttpClient) SetRetryCount(retry int) *HttpClient {
//	client := *c
//	client.retryCount = retry
//	return &client
//}
//
//func (c *HttpClient) SetRetryWaitTime(waitTime time.Duration) *HttpClient {
//	client := *c
//	client.retryWaitTime = waitTime
//	return &client
//}

func (c *HttpClient) SetDialTimeout(timeout time.Duration) *HttpClient {
	client := *c
	client.Transport = &http.Transport{
		DialContext: (&net.Dialer{Timeout: timeout}).DialContext,
		TLSHandshakeTimeout: timeout,
	}
	return &client
}

func (c *HttpClient) SetRequestTimeout(timeout time.Duration) *HttpClient {
	client := *c
	client.Timeout = timeout
	return &client
}

func (c *HttpClient) Set(key string, value string) *HttpClient {
	client := *c
	client.headers.Set(key, value)
	return &client
}

// SetHeader will set the header and replace any existing header
func (c *HttpClient) SetHeader(header http.Header) *HttpClient {
	client := *c
	client.headers = header
	return &client
}

func (c *HttpClient) Body(requestBody []byte) *HttpClient {
	client := *c
	client.requestBody = requestBody
	return &client
}

func (c *HttpClient) BodyWithType(requestBody []byte, contentType HttpCallContentType) *HttpClient {
	client := *c
	client.requestBody = requestBody
	client.headers.Set("Content-Type", string(contentType))
	return &client
}

func (c *HttpClient) AddFormData(key string, values ...string) *HttpClient {
	client := *c
	if client.form == nil {
		client.form = make(url.Values)
	}
	client.form[key] = values
	return &client
}

func (c *HttpClient) Call(options ...Option) (*http.Response, error) {
	clopts := &clientOptions{
		RetryCount:    0,
		RetryWaitTime: 0,
		KeepAlive:     true,
	}

	for _, o := range options {
		o(clopts)
	}

	body := bytes.NewReader(c.requestBody)
	if len(c.form) > 0 {
		body = bytes.NewReader([]byte(c.form.Encode()))
		c.headers.Set("Content-Type", string(HTTP_CALL_CONTENT_URLFORM))
	}

	req, err := http.NewRequest(c.method, c.url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to form http request: %s", err)
	}

	req.Header = c.headers
	req.Close = !clopts.KeepAlive

	var res *http.Response
	if clopts.RetryCount == 0 {
		return c.execute(req)
	}

	for i := 0; i < clopts.RetryCount; i++ {
		res, err = c.execute(req)
		if err == nil {
			if res.StatusCode != http.StatusTooManyRequests {
				break
			}
		}

		waitTime := clopts.RetryWaitTime
		if waitTime == 0 {
			waitTime = defaultWaitTime
		}

		time.Sleep(waitTime)
	}

	return res, err
}

func (c *HttpClient) execute(req *http.Request) (*http.Response, error) {
	if len(c.form) > 0 {
		if c.method != "POST" {
			return nil, fmt.Errorf("cannot send Form %q with '%s' method", c.form.Encode(), c.method)
		}
	}

	return c.Do(req)
}

func (c *HttpClient) DumpRequest() ([]byte, error) {
	req, err := http.NewRequest(c.method, c.url, bytes.NewBuffer(c.requestBody))
	if err != nil {
		return nil, err
	}

	return httputil.DumpRequest(req, true)
}

func Retry(count int, waitTime time.Duration) Option {
	return func(options *clientOptions) {
		options.RetryCount = count
		options.RetryWaitTime = waitTime
	}
}

func DisableKeepAlive() Option {
	return func(options *clientOptions) {
		options.KeepAlive = false
	}
}
