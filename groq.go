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

const DEFAULT_MODEL = "llama-3.1-8b-instant"
const DEFAULT_TTS_MODEL = "playai-tts"
const DEFAULT_TTS_VOICE = "Fritz-PlayAI"
const DEFAULT_TTS_FORMAT = "wav"

func (client *GroqClient) buildQueryRequest(query string, params QueryParameters) (ChatCompletionRequest, error) {
	if client.ApiKey == "" {
		return ChatCompletionRequest{}, errors.New("API key is not specified")
	}
	if client.Model == "" {
		fmt.Println("Model is not specified. Defaulting to `" + DEFAULT_MODEL + "`")
		client.Model = DEFAULT_MODEL
	}
	req := ChatCompletionRequest{
		Model:    client.Model,
		Messages: []Message{},
	}
	if params.SystemPrompt != "" {
		req.Messages = append(req.Messages, Message{Role: "system", Content: params.SystemPrompt})
	}
	req.Messages = append(req.Messages, Message{Role: "user", Content: query})
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

func (client *GroqClient) buildSpeechRequest(text string, params SpeechParameters) (SpeechRequest, error) {
    if client.ApiKey == "" {
        return SpeechRequest{}, errors.New("API key is not specified")
    }
    model := DEFAULT_TTS_MODEL
    voice := DEFAULT_TTS_VOICE
    responseFormat := DEFAULT_TTS_FORMAT
    if params.Voice != "" {
        voice = params.Voice
    }
    if params.ResponseFormat != "" {
        responseFormat = params.ResponseFormat
    }
    return SpeechRequest{
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

func (client *GroqClient) GetRequestsCount() int {
	return client.RequestsCount
}
