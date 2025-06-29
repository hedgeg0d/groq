package groq

type QueryParameters struct {
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	TopP         float64 `json:"top_p,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Delta struct {
	Content string `json:"content"`
}

type StreamChoice struct {
	Delta Delta `json:"delta"`
}

type ChatCompletionStreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}

type SpeechParameters struct {
    Voice          string `json:"voice,omitempty"`
    ResponseFormat string `json:"response_format,omitempty"`
}

type SpeechRequest struct {
    Model          string `json:"model"`
    Input          string `json:"input"`
    Voice          string `json:"voice"`
    ResponseFormat string `json:"response_format"`
}

type SpeechResponse struct {
    Audio []byte `json:"audio"`
}