package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq"
)


func main() {
	client := groq.GroqClient{
		ApiKey: os.Getenv("GROQ_API_KEY"),
		Model: "llama-3.1-8b-instant",
	}
	resp, err := client.Ask("Hello. Tell me about yourself")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Output: " + resp)
	}
}