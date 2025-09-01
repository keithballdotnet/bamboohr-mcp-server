package main

import (
	"testing"
	"time"
)

func TestNewBambooHRClient(t *testing.T) {
	client := NewBambooHRClient("testcompany", "testkey")

	if client.Company != "testcompany" {
		t.Errorf("Expected company to be 'testcompany', got '%s'", client.Company)
	}

	if client.APIKey != "testkey" {
		t.Errorf("Expected API key to be 'testkey', got '%s'", client.APIKey)
	}

	expectedURL := "https://testcompany.bamboohr.com/api/gateway.php/testcompany/v1"
	if client.BaseURL != expectedURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", expectedURL, client.BaseURL)
	}

	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30 seconds, got %v", client.HTTPClient.Timeout)
	}
}

func TestBambooHRClientValidation(t *testing.T) {
	tests := []struct {
		name     string
		company  string
		apiKey   string
		hasError bool
	}{
		{"Valid inputs", "mycompany", "sk_12345", false},
		{"Empty company", "", "sk_12345", false},  // Client creation doesn't validate, main() does
		{"Empty API key", "mycompany", "", false}, // Client creation doesn't validate, main() does
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewBambooHRClient(tt.company, tt.apiKey)
			if client == nil {
				t.Error("Expected client to be created, got nil")
			}
		})
	}
}
