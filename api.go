package groq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func createChatCompletion(req chatCompletionRequest, apiKey string) (*chatCompletionResponse, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var chatResp chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}
	return &chatResp, nil
}

func createChatCompletionStream(req chatCompletionRequest, apiKey string) (chan string, error) {
	req.Stream = true
	url := "https://api.groq.com/openai/v1/chat/completions"
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	chunks := make(chan string)
	go func() {
		defer resp.Body.Close()
		defer close(chunks)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "data: [DONE]" {
				break
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			var streamResp struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}

			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				chunks <- streamResp.Choices[0].Delta.Content
			}
		}
	}()

	return chunks, nil
}

func createSpeech(req speechRequest, apiKey string) ([]byte, error) {
	url := "https://api.groq.com/openai/v1/audio/speech"
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): failed to parse error response", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error (status %d): %s (%s)", resp.StatusCode, errResp.Error.Message, errResp.Error.Type)
	}
	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return audio, nil
}

func createTranscription(apiKey string, audioData []byte, model string, params TranscriptionParameters) (*transcriptionResponse, error) {
	url := "https://api.groq.com/openai/v1/audio/transcriptions"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, bytes.NewReader(audioData))
	if err != nil {
		return nil, fmt.Errorf("failed to copy audio data to form: %v", err)
	}

	if err := writer.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("failed to write model field: %v", err)
	}

	if params.Language != "" {
		if err := writer.WriteField("language", params.Language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %v", err)
		}
	}
	if params.Prompt != "" {
		if err := writer.WriteField("prompt", params.Prompt); err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %v", err)
		}
	}

	if params.Temperature != 0.0 {
		tempStr := fmt.Sprintf("%f", params.Temperature)
		if err := writer.WriteField("temperature", tempStr); err != nil {
			return nil, fmt.Errorf("failed to write temperature field: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): failed to parse error response", resp.StatusCode)
		}
		return nil, fmt.Errorf("API error (status %d): %s (%s)", resp.StatusCode, errResp.Error.Message, errResp.Error.Type)
	}

	var transResp transcriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transResp); err != nil {
		return nil, fmt.Errorf("failed to decode transcription response: %v", err)
	}

	return &transResp, nil
}
