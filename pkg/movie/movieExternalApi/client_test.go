package movieExternalApi

import (
	"testing"
	"time"
)

func TestNewAPIClient(t *testing.T) {
	baseURL := "http://testapi.com"
	client := NewAPIClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %q, got %q", baseURL, client.BaseURL)
	}
	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}
	if client.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("Expected HTTPClient timeout of 10s, got %v", client.HTTPClient.Timeout)
	}
}
