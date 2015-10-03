package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestShouldBeAbleToSetQuantityOfRetries(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	if req, _ := NewRequest("GET", ts.URL, nil); req.retry != Retries {
		t.Error(fmt.Sprintf("retry should've been set to %d, got %d", Retries, req.retry))
	}
}

func TestShouldBeAbleToSetURL(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	if req, _ := NewRequest("GET", ts.URL, nil); req.URL.String() != ts.URL {
		t.Error(fmt.Sprintf("url should've been set to %s, got %s", ts.URL, req.URL.String()))
	}
}

func TestShouldRetryAfterResponseCode5XX(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Do()

	if req.retry != 0 {
		t.Error(fmt.Sprintf("retry should be 0, got %d", req.retry))
	}
}

func TestShouldRetryWhenErrorHappens(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errors.New("emit macho dwarf: elf header corrupted")
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", "http://www.qoroqer.com", nil)
	req.Do()

	if req.retry != 0 {
		t.Error(fmt.Sprintf("retry should be 0, got %d", req.retry))
	}
}

func TestShouldKeepRetryCountIntactWhenOK(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			js, _ := json.Marshal(&numbers)
			w.Write(js)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Do()

	if req.retry != Retries {
		t.Error(fmt.Sprintf("retry should not have changed from %d, to %d", Retries, req.retry))
	}
}

func TestShouldRetrieWhenTimeout(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Timeout(1 * time.Nanosecond)

	req.Do()

	if req.retry != 0 {
		t.Error(fmt.Sprintf("retry should be 0, got %d", req.retry))
	}
}

func TestResponseShouldNotBeNilWhenOK(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			js, _ := json.Marshal(&numbers)
			w.Write(js)
		}),
	)
	defer ts.Close()

	req, _ := NewRequest("GET", ts.URL, nil)
	req.Do()

	if resp, _ := req.Do(); resp == nil {
		t.Error("response should be nil")
	}

}
