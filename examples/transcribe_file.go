package main

import (
	"fmt"
	"os"
	"github.com/hedgeg0d/groq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run examples/transcribe_file.go <path_to_audio_file>")
		return
	}
	audioData, _ := os.ReadFile(os.Args[1])
	client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
	text, _ := client.CreateTranscription(audioData, groq.TranscriptionParameters{})
	fmt.Println("Transcription: " + text)
}
