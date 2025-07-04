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

const default_model = "llama-3.1-8b-instant"
const default_tts_model = "playai-tts"
const default_tts_voice = "Fritz-PlayAI"
const default_tts_format = "wav"
const default_transcription_model = "whisper-large-v3"

func (client *GroqClient) buildQueryRequest(query string, params QueryParameters) (chatCompletionRequest, error) {
	if client.ApiKey == "" {
		return chatCompletionRequest{}, errors.New("API key is not specified")
	}
	if client.Model == "" {
		fmt.Println("Model is not specified. Defaulting to `" + default_model + "`")
		client.Model = default_model
	}
	req := chatCompletionRequest{
		Model:    client.Model,
		Messages: []message{},
	}
	if params.SystemPrompt != "" {
		req.Messages = append(req.Messages, message{Role: "system", Content: params.SystemPrompt})
	}
	req.Messages = append(req.Messages, message{Role: "user", Content: query})
	if params.MaxTokens != 0 {
		req.MaxTokens = params.MaxTokens
	}
	if params.Temperature != 0.0 {
		req.Temperature = params.Temperature
	}
	if params.TopP != 0 {
		req.TopP = params.TopP
	}
	return req, nil
}

func (client *GroqClient) buildSpeechRequest(text string, params SpeechParameters) (speechRequest, error) {
	if client.ApiKey == "" {
		return speechRequest{}, errors.New("API key is not specified")
	}
	model := default_tts_model
	voice := default_tts_voice
	responseFormat := default_tts_format
	if params.Voice != "" {
		voice = params.Voice
	}
	if params.ResponseFormat != "" {
		responseFormat = params.ResponseFormat
	}
	return speechRequest{
		Model:          model,
		Input:          text,
		Voice:          voice,
		ResponseFormat: responseFormat,
	}, nil
}

func (client *GroqClient) Ask(query string) (string, error) {
	return client.Query(query, QueryParameters{})
}

func (client *GroqClient) Query(query string, params QueryParameters) (string, error) {
	req, err := client.buildQueryRequest(query, params)
	if err != nil {
		return "", err
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
	req, err := client.buildQueryRequest(query, params)
	if err != nil {
		return nil, err
	}
	req.Stream = true
	chunks, err := createChatCompletionStream(req, client.ApiKey)
	if err != nil {
		return nil, err
	}
	client.RequestsCount++
	return chunks, nil
}

func (client *GroqClient) CreateSpeech(text string, params SpeechParameters) ([]byte, error) {
	req, err := client.buildSpeechRequest(text, params)
	if err != nil {
		return nil, err
	}
	audio, err := createSpeech(req, client.ApiKey)
	if err != nil {
		return nil, err
	}
	client.RequestsCount++
	return audio, nil
}

func (client *GroqClient) CreateTranscription(audioData []byte, params TranscriptionParameters) (string, error) {
	if client.ApiKey == "" {
		return "", errors.New("API key is not specified")
	}

	model := default_transcription_model
	resp, err := createTranscription(client.ApiKey, audioData, model, params)
	if err != nil {
		return "", err
	}

	client.RequestsCount++
	return resp.Text, nil
}

func (client *GroqClient) GetRequestsCount() int {
	return client.RequestsCount
}
