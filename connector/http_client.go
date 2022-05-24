package connector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/NuttapolCha/test-band-data-feeder/log"

	"net/http"
)

type CustomHttpClient struct {
	logger log.Logger
}

func NewCustomHttpClient(logger log.Logger) *CustomHttpClient {
	return &CustomHttpClient{
		logger: logger,
	}
}

func (c *CustomHttpClient) Get(endpoint string, queryStr map[string]string, retryCount int) ([]byte, error) {
	logger := c.logger

	start := time.Now()
	var err error

	// establish a new request
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Errorf("could not establish a new request to %s because: %v", endpoint, err)
		return nil, err
	}

	// construct query strings
	q := req.URL.Query()
	for key, val := range queryStr {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	// requests up to defined attempts
	client := http.DefaultClient
	for attmps := 0; attmps < retryCount+1; attmps++ {
		var resp *http.Response
		var respBody []byte

		logger.Debugf("attemp: %d requesting to %s", attmps+1, endpoint)
		resp, err = client.Do(req)
		if err != nil {
			logger.Errorf("attempt: %d could not GET Request to %s because: %v, will retry in 1 seconds", attmps+1, endpoint, err)
			time.Sleep(1 * time.Second)
			continue
		}

		respBody, err = c.resolveRespResult(resp)
		if err != nil {
			logger.Errorf("attempt: %d could not resolve response result from %s because: %v, will retry in 1 seconds", attmps+1, endpoint, err)
			time.Sleep(1 * time.Second)
			continue
		}

		logger.Debugf("request to %s success, time used: %v", endpoint, time.Since(start))
		return respBody, nil
	}

	// attemps have been reached the maximum
	logger.Errorf("could not GET Request to %s after %d attempts because: %v", endpoint, retryCount, err)
	return nil, err
}

func (c *CustomHttpClient) PostJSON(endpoint string, body []byte, retryCount int) ([]byte, error) {
	logger := c.logger

	start := time.Now()
	var err error

	// requests up to defined attempts
	for attmps := 0; attmps < retryCount+1; attmps++ {
		var resp *http.Response
		var respBody []byte

		logger.Debugf("attempt: %d requesting to %s", attmps+1, endpoint)
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(body))
		if err != nil {
			logger.Errorf("attempt: %d could not POST Request to %s because: %v will retry in 1 seconds", attmps+1, endpoint, err)
			time.Sleep(1 * time.Second)
			continue
		}

		respBody, err = c.resolveRespResult(resp)
		if err != nil {
			logger.Errorf("attempt: %d could not resolve response result from %s because: %v, will retry in 1 seconds", attmps+1, endpoint, err)
			time.Sleep(1 * time.Second)
			continue
		}

		logger.Debugf("request to %s success, time used: %v", endpoint, time.Since(start))
		return respBody, nil
	}

	// attemps have been reached the maximum
	logger.Errorf("could not POST Request to %s after %d attempts because: %v", endpoint, retryCount, err)
	return nil, err
}

func (c *CustomHttpClient) resolveRespResult(resp *http.Response) ([]byte, error) {
	logger := c.logger

	var err error

	// check resposne status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("got unexpected response %s", resp.Status)
	}

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("response body =")
	logger.BeautyJSON(body)
	return body, nil
}
