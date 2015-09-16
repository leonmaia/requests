package requests

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cenkalti/backoff"
)

func doReq(r *Request) (*http.Response, error) {
	res, err := http.Get(r.URL.String())
	if err != nil && r.retry > 0 {
		r.retry--
		return nil, err
	}
	if res != nil && res.StatusCode >= 500 && res.StatusCode <= 599 && r.retry > 0 {
		r.retry--
		return nil, errors.New("Server Error")
	}

	return res, nil
}

// Do should be called when the Request is fully configured.
func (r *Request) Do() ([]byte, error) {
	res, err := doReq(r)
	if err != nil {
		op := r.operation()
		err = backoff.Retry(op, r.backoff)
		if err != nil {
			return nil, err
		}
	}

	if res != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	return nil, errors.New("Server Error")
}

func (r *Request) operation() func() error {
	return func() error {
		_, err := doReq(r)
		return err
	}
}
