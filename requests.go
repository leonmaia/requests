// Package requests implements functions to manipulate requests.
package requests

import (
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
)

const (
	// Default quantity of retries
	Retries = 3
	// Default timeout is 30 seconds
	Timeout = 30 * time.Second
)

// NewRequest returns an Request with exponential backoff as default.
func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	return &Request{
			req, Retries, Timeout, backoff.NewExponentialBackOff()},
		nil
}

// Request type.
type Request struct {
	*http.Request
	retry   int
	timeout time.Duration
	backoff *backoff.ExponentialBackOff // Default Type of backoff.
}

// Set the amount of retries
func (r *Request) Retries(times int) *Request {
	r.retry = times
	return r
}

// Timeout specifies a time limit for requests made by the Client.
// A Timeout of zero means no timeout.
func (r *Request) Timeout(t time.Duration) *Request {
	r.timeout = t
	return r
}

// New Client with timeout
func (r *Request) newClient() *http.Client {
	return &http.Client{Timeout: r.timeout}
}
