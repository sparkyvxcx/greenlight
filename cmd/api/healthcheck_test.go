package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.sparkyvxcx.co/internal/assert"
	"greenlight.sparkyvxcx.co/internal/jsonlog"
)

func BenchmarkEncoder(b *testing.B) {
	w := httptest.NewRecorder()
	r := new(http.Request)
	for n := 0; n < b.N; n++ {
		healthcheckHandlerWithEncoder(w, r)
	}
}

func BenchmarkMarshal(b *testing.B) {
	w := httptest.NewRecorder()
	r := new(http.Request)
	for n := 0; n < b.N; n++ {
		healthcheckHandlerWithMarshal(w, r, false)
	}
}

func BenchmarkMarshalIndent(b *testing.B) {
	w := httptest.NewRecorder()
	r := new(http.Request)
	for n := 0; n < b.N; n++ {
		healthcheckHandlerWithMarshal(w, r, true)
	}
}

func TestHealthCheck(t *testing.T) {
	logger := jsonlog.New(io.Discard, jsonlog.LevelInfo)
	app := &application{
		logger: logger,
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/v1/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rs.Header.Get("Content-Type"), "application/json")

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	responseBody := string(body)
	assert.Contains(t, responseBody, "status")
	assert.Contains(t, responseBody, "system_info")
	assert.Contains(t, responseBody, "environment")
	assert.Contains(t, responseBody, "version")
}
