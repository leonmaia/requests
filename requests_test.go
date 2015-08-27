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

	req := New()
	req.URL = ts.URL
	req.Retry = retry

	c.Assert(req.Retry, Equals, retry)
	c.Assert(req.URL, Equals, ts.URL)

	_, err := req.Do()

	if err != nil {
		fmt.Println("err", err)
	}

	c.Assert(req.Retry, Equals, 0)
}

func (s *TestSuite) TestRetriesWithError(c *C) {
	retry := 3
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errors.New("emit macho dwarf: elf header corrupted")
		}),
	)
	defer ts.Close()

	req := New()
	req.URL = "qojerq"
	req.Retry = retry

	c.Assert(req.Retry, Equals, retry)
	c.Assert(req.URL, Equals, "qojerq")

	_, err := req.Do()

	if err != nil {
		fmt.Println("err", err)
	}

	c.Assert(req.Retry, Equals, 0)
}
