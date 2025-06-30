package main

import (
	"fmt"
	"os"
	"github.com/hedgeg0d/groq"
)

func main() {
	client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
	audioData, _ := client.CreateSpeech("The quick brown fox jumps over the lazy dog.", groq.SpeechParameters{})
	params := groq.TranscriptionParameters{
		Language:    "en",
		Prompt:      "This is a test.",
		Temperature: 0.2,
	}
	text, _ := client.CreateTranscription(audioData, params)
	fmt.Println("Transcription" + text)
}
