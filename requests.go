// Package requests implements functions to manipulate requests.
package requests

import (
	"io"
	"net/http"

	"github.com/cenkalti/backoff"
)

// Default quantity of retries
const RETRIES = 3

// New returns an Request with exponential backoff as default.
func NewRequest(method, urlStr string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	return &Request{req, RETRIES, backoff.NewExponentialBackOff()}, nil
}

// Request type.
type Request struct {
	*http.Request
	retry   int
	backoff *backoff.ExponentialBackOff // Default Type of backoff.
}

// Set the amount of retries
func (r *Request) Retries(times int) *Request {
	r.retry = times
	return r
}
