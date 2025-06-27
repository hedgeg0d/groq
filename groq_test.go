package groq

import (
	"github.com/hedgeg0d/groq"
	"testing"
)

func TestGroqClient_Query_Validation(t *testing.T) {
	client := &groq.GroqClient{ApiKey: "", Model: ""}
	_, err := client.Query("Test prompt", groq.QueryParameters{})
	if err == nil || err.Error() != "API key is not specified" {
		t.Errorf("Expected error 'API key is not specified', got %v", err)
	}

	client.ApiKey = "test-api-key"
	client.Model = ""
	_, err = client.Query("Test prompt", groq.QueryParameters{MaxTokens: 100, Temperature: 0.7, TopP: 0.9})
	if err.Error() != "unexpected status code: 401" {
		t.Errorf("Expected \"unexpected status code: 401\", got %v", err)
	}
	if client.Model != "llama-3.1-8b-instant" {
		t.Errorf("Expected client.Model 'llama-3.1-8b-instant', got %s", client.Model)
	}
	if client.RequestsCount != 0 {
		t.Errorf("Expected RequestsCount 0, got %d", client.RequestsCount)
	}
}
