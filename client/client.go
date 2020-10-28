package client

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client    *http.Client
	delay     time.Duration
	stopCh    chan struct{}
	requestCh chan struct{}
}

func NewHttpClient(timeout time.Duration, delay time.Duration) *Client {
	client := new(Client)
	client.delay = delay
	client.stopCh = make(chan struct{})
	client.requestCh = make(chan struct{})
	client.client = &http.Client{
		Timeout: timeout,
	}
	return client
}

func (c *Client) freeRequestCh() {
	select {
	case <-c.requestCh:
	default:
	}
}

func (c *Client) delayerProcess() {
	if c.IsClosed() {
		return
	}
	timer := time.NewTimer(c.delay)
	defer timer.Stop()
	for {
		if c.IsClosed() {
			return
		}
		select {
		case <-c.stopCh:
			return
		case <-timer.C:
			c.freeRequestCh()
			timer.Reset(c.delay)
		}
	}
}

func (c *Client) StartDelayer() error {
	if c.IsClosed() {
		return fmt.Errorf("can't start delayer client stopped")
	}
	go c.delayerProcess()
	return nil
}

func (c *Client) IsClosed() bool {
	select {
	case <-c.stopCh:
		return true
	default:
		return false
	}
}

func (c *Client) Close() {
	if !c.IsClosed() {
		close(c.stopCh)
	}
}

func (c *Client) doSimpleRequest(method, urlPath string, urlParams, headers map[string]string, body []byte, response interface{}) error {
	if urlParams != nil {
		params := url.Values{}
		for key, val := range headers {
			params.Add(key, val)
		}
		encodedParams := params.Encode()
		if len(encodedParams) > 0 {
			urlPath = urlPath + "?" + encodedParams
		}
	}
	req, err := http.NewRequest(method, urlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		if err != nil {
			return err
		}
		err = jsoniter.Unmarshal(bodyBytes, response)
		if err != nil {
			return fmt.Errorf("can't unmarshal response; %s", string(bodyBytes))
		}
	} else {
		return fmt.Errorf("response status is not OK; %s", string(bodyBytes))
	}
	return nil
}

func (c *Client) doRequestWithDelay(method, url string, urlParams, headers map[string]string, body []byte, response interface{}) error {
	if c.IsClosed() {
		return fmt.Errorf("client is closed")
	}
	c.requestCh <- struct{}{}
	err := c.doSimpleRequest(method, url, urlParams, headers, body, response)

	return err
}

func (c *Client) DoRequest(method, url string, urlParams, headers map[string]string, body []byte, response interface{}) error {
	return c.doRequestWithDelay(method, url, urlParams, headers, body, response)
}
