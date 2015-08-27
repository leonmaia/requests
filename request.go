package request

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/cenkalti/backoff"
)

type Request struct {
	URL     string
	Retry   int
	body    []byte
	backoff *backoff.ExponentialBackOff
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

func (r *Request) do() ([]byte, error) {
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

func NewRequest() *Request {
	r := &Request{}
	r.backoff = backoff.NewExponentialBackOff()
	return r
}
