package groq

import (
	"errors"
	"fmt"
)

type GroqClient struct {
	ApiKey        string
	Model         string
	RequestsCount int
}

type QueryParameters struct {
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

func (client *GroqClient) Ask(query string) (string, error) {
	if client.ApiKey == "" {
		return "", errors.New("API key is not specified")
	}
	if client.Model == "" {
		fmt.Println("Model is not specified. Defaulting to `llama-3.1-8b-instant`")
		client.Model = "llama-3.1-8b-instant"
	}
	req := ChatCompletionRequest{
		Model: client.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	resp, err := createChatCompletion(req, client.ApiKey)
	if err != nil {
		return "", err
	} else {
		client.RequestsCount++
		return resp.Choices[0].Message.Content, nil
	}
}

func (client *GroqClient) Query(query string, params QueryParameters) (string, error) {
	if client.ApiKey == "" {
		return "", errors.New("API key is not specified")
	}
	if client.Model == "" {
		fmt.Println("Model is not specified. Defaulting to `llama-3.1-8b-instant`")
		client.Model = "llama-3.1-8b-instant"
	}

	req := ChatCompletionRequest{
		Model: client.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
	}
	if params.MaxTokens != 0 {
		req.MaxTokens = params.MaxTokens
	}
	if params.Temperature != 0.0 {
		req.Temperature = params.Temperature
	}
	if params.TopP != 0 {
		req.TopP = params.TopP
	}

	resp, err := createChatCompletion(req, client.ApiKey)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices received")
	}

	client.RequestsCount++
	return resp.Choices[0].Message.Content, nil
}

func (client *GroqClient) AskQueryStream(query string, params QueryParameters) (chan string, error) {
	if client.ApiKey == "" {
		return nil, fmt.Errorf("ApiKey is not specified")
	}
	if client.Model == "" {
		fmt.Println("Model is not specified. Defaulting to `llama-3.1-8b-instant`")
		client.Model = "llama-3.1-8b-instant"
	}

	req := ChatCompletionRequest{
		Model: client.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Stream: true,
	}
	if params.MaxTokens != 0 {
		req.MaxTokens = params.MaxTokens
	}
	if params.Temperature != 0.0 {
		req.Temperature = params.Temperature
	}
	if params.TopP != 0 {
		req.TopP = params.TopP
	}

	chunks, err := createChatCompletionStream(req, client.ApiKey)
	if err != nil {
		return nil, err
	}

	client.RequestsCount++
	return chunks, nil
}
