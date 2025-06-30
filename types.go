package groq

type QueryParameters struct {
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	TopP         float64 `json:"top_p,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type choice struct {
	Message message `json:"message"`
}

type chatCompletionResponse struct {
	Choices []choice `json:"choices"`
}

type delta struct {
	Content string `json:"content"`
}

type streamChoice struct {
	Delta delta `json:"delta"`
}

type chatCompletionStreamResponse struct {
	Choices []streamChoice `json:"choices"`
}

type SpeechParameters struct {
	Voice          string `json:"voice,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

type speechRequest struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Voice          string `json:"voice"`
	ResponseFormat string `json:"response_format"`
}

type speechResponse struct {
	Audio []byte `json:"audio"`
}

type TranscriptionParameters struct {
	Language    string  `json:"language,omitempty"`
	Prompt      string  `json:"prompt,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type transcriptionResponse struct {
	Text string `json:"text"`
}
