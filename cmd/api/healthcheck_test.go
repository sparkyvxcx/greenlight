package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
		healthcheckHandlerWithMarshal(w, r)
	}
}
