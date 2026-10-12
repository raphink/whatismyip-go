package whatismyip

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
	handler := http.HandlerFunc(whatIsMyIP)

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

func TestGetAllowedIPsWithoutPlayID(t *testing.T) {
	originalProjectID, originalClient := projectID, client
	projectID, client = "test-project", nil
	defer func() {
		projectID, client = originalProjectID, originalClient
	}()

	allowedIPs, err := getAllowedIPs("")
	if err != nil {
		t.Fatalf("getAllowedIPs returned an error: %v", err)
	}
	if len(allowedIPs) != 0 {
		t.Fatalf("getAllowedIPs returned %v, want no allowed IPs", allowedIPs)
	}
}
