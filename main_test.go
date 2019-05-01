package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httphandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code got: %v want %v",
			status, http.StatusOK)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("Unexpected body got: %v want %v",
			rr.Body.String(), expected)
	}
}
