package requests

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cenkalti/backoff"
)

func doReq(r *Request, c *http.Client) (*http.Response, error) {
	res, err := c.Do(r.Request)
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

func ToBytes(res *http.Response) ([]byte, error) {
	if res != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}
	return nil, errors.New("Invalid response")
}

func ToByteChannel(res *http.Response, bufferSize int) <-chan []byte {
	resp_chan := make(chan []byte)
	b := bufio.NewReader(res.Body)
	buf := make([]byte, 0, bufferSize)
	go func() {
		for {
			n, err := b.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				if err == nil {
					continue
				}
				if err == io.EOF {
					close(resp_chan)
					break
				}
				log.Fatal(err)
			}
			resp_chan <- buf
		}
	}()
	return resp_chan
}

func ToJsonChannel(r *http.Response, dec interface{}) <-chan interface{} {
	resp_chan := make(chan interface{})
	decoder := json.NewDecoder(r.Body)
	go func() {
		for {
			err := decoder.Decode(dec)
			if err == nil {
				resp_chan <- dec
			} else if err == io.EOF {
				defer close(resp_chan)
				break
			} else if err != nil {
				fmt.Println(err)
				defer close(resp_chan)
				break
			}
		}
	}()

	return resp_chan
}

// Do should be called when the Request is fully configured.
func (r *Request) Do() (*http.Response, error) {
	c := r.newClient()
	res, err := doReq(r, c)
	if err != nil {
		op := r.operation(c)
		err = backoff.Retry(op, r.backoff)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (r *Request) operation(c *http.Client) func() error {
	return func() error {
		_, err := doReq(r, c)
		return err
	}
}
