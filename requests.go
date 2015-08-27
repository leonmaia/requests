// Package requests implements functions to manipulate requests.
package requests

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cenkalti/backoff"
)

// New returns an Request with exponential backoff as default.
func New() *Request {
	return &Request{backoff: backoff.NewExponentialBackOff()}
}

// Request type.
type Request struct {
	URL     string
	Retry   int                         // Amount of retries.
	body    []byte                      // Response Body.
	backoff *backoff.ExponentialBackOff // Default Type of backoff.
}

func doReq(r *Request) error {
	res, err := http.Get(r.URL)
	if err != nil && r.Retry > 0 {
		r.Retry--
		return err
	}
	if res != nil && res.StatusCode >= 500 && res.StatusCode <= 599 && r.Retry > 0 {
		r.Retry--
		return errors.New("Server Error")
	}

	if res != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		r.body = body
	}

	return nil
}

// Do should be callend when the Request is fully configured.
func (r *Request) Do() ([]byte, error) {
	err := doReq(r)
	if err != nil {
		op := r.operation()
		err = backoff.Retry(op, r.backoff)
		if err != nil {
			return nil, err
		}
	}

	return r.body, nil
}

func (r *Request) operation() func() error {
	return func() error {
		return doReq(r)
	}
}
