package helpers

import (
	"io"
	"net/http"
)

type RetryClient struct {
	http.Client
}

func NewRetryClient() *RetryClient {
	return &RetryClient{}
}

func (rc *RetryClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	attempt := 0
	for {
		if attempt > 3 {
			return
		}
		attempt++
		resp, err = rc.Client.Post(url, contentType, body)
		if err != nil {
			return
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 420 {
			continue
		}
		return
	}
}
