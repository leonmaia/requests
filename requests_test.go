package requests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestShouldRetryAfterResponseCode5XX(c *C) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)

	c.Assert(req.retry, Equals, Retries)
	c.Assert(req.httpReq.URL.String(), Equals, ts.URL)

	req.Do()

	c.Assert(req.retry, Equals, 0)
}

func (s *TestSuite) TestShouldRetryWhenErrorHappens(c *C) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errors.New("emit macho dwarf: elf header corrupted")
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", "http://www.qoroqer.com", nil)

	c.Assert(req.httpReq.URL.String(), Equals, "http://www.qoroqer.com")

	req.Do()

	c.Assert(req.retry, Equals, 0)
}

func (s *TestSuite) TestShouldKeepRetryCountIntactWhenOK(c *C) {
	numbers := []int{1, 2, 3, 4, 5}
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			js, _ := json.Marshal(&numbers)
			w.Write(js)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)

	c.Assert(req.retry, Equals, Retries)
	c.Assert(req.httpReq.URL.String(), Equals, ts.URL)

	resp, _ := req.Do()

	c.Assert(req.retry, Equals, Retries)
	c.Assert(resp, Not(Equals), nil)
}

func (s *TestSuite) TestShouldRetrieWhenTimeout(c *C) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Nanosecond)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Timeout(1 * time.Nanosecond)

	c.Assert(req.retry, Equals, Retries)
	c.Assert(req.httpReq.URL.String(), Equals, ts.URL)

	req.Do()

	c.Assert(req.retry, Equals, 0)
}

func (s *TestSuite) TestShouldChangeQuantityOfRetries(c *C) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Retries(2)

	c.Assert(req.retry, Equals, 2)
	c.Assert(req.httpReq.URL.String(), Equals, ts.URL)

	req.Do()

	c.Assert(req.retry, Equals, 0)
}
