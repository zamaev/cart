package middleware

import (
	"context"
	"io"
	"net/http"
	"route256/cart/pkg/tracing"
	"strconv"

	"go.opentelemetry.io/otel/codes"
)

type RetryClient struct {
	http.Client
}

func NewRetryClient() *RetryClient {
	return &RetryClient{}
}

func (rc *RetryClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	_, span := tracing.Start(ctx, "RetryClient.Post")
	defer span.End()
	defer func() {
		span.AddEvent("StatusCode: " + strconv.Itoa(resp.StatusCode))
	}()

	attempt := 0
	for {
		if attempt > 3 {
			return
		}
		attempt++
		span.AddEvent("Attempt: " + strconv.Itoa(attempt))
		resp, err = rc.Client.Post(url, contentType, body)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 420 {
			continue
		}
		return
	}
}
