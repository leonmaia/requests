package requests

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestRetriesWith500(c *C) {
	retry := 3
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Retries(retry)

	c.Assert(req.retry, Equals, retry)
	c.Assert(fmt.Sprintf("%s://%s", req.URL.Scheme, req.URL.Host), Equals, ts.URL)

	_, err := req.Do()

	if err != nil {
		fmt.Println("err", err)
	}

	c.Assert(req.retry, Equals, 0)
}

func (s *TestSuite) TestRetriesWithError(c *C) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errors.New("emit macho dwarf: elf header corrupted")
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", "http://www.qoroqer.com", nil)

	c.Assert(req.URL.Host, Equals, "www.qoroqer.com")

	_, err := req.Do()

	if err != nil {
		fmt.Println("err", err)
	}

	c.Assert(req.retry, Equals, 0)
}
