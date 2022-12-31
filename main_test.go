package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	streams = make(map[int]string)
	heartbeats = make(map[int]time.Time)
	key = "secretkey"

	code := m.Run()

	os.Exit(code)
}

func TestWhepEndpoint(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/whep/endpoint/1234", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(whepEndpoint)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	rr = httptest.NewRecorder()
	streams[1234] = "http://foobar/1234"
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	redirect := rr.Result().Header.Get("Location")

	if redirect != streams[1234] {
		t.Errorf("handler returned unexpected redirect: got %v want %v",
			redirect, streams[1234])
	}
}

func TestStartStream(t *testing.T) {
	form := url.Values{}
	form.Add("channel_id", "12345")
	form.Add("endpoint", "http://foobar/12345")
	req, err := http.NewRequest("POST", "/v1/state/start_stream", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", key)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(startStream)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	endpoint, ok := streams[12345]
	if !ok || endpoint != "http://foobar/12345" {
		t.Error("stream was not added to state successfully")
	}
}

func TestEndStream(t *testing.T) {
	form := url.Values{}
	form.Add("channel_id", "123456")
	req, err := http.NewRequest("POST", "/v1/state/end_stream", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", key)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(endStream)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	endpoint, ok := streams[123456]
	if ok || endpoint == "http://foobar/123456" {
		t.Error("stream was not removed from state successfully")
	}
}

func TestHeartbeatTimeout(t *testing.T) {
	streams[4321] = "http://foobar/4321"
	heartbeats[4321] = time.Now().Add(time.Second * -10)

	checkForDeadChannels(time.Duration(time.Second * 5))

	_, ok := streams[4321]
	if ok {
		t.Error("stream was not removed from state successfully")
	}
}
