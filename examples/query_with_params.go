package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq"
)


func main() {
	client := groq.GroqClient{
		ApiKey: os.Getenv("GROQ_API_KEY"),
		Model: "",
	}
	params := groq.QueryParameters {
		MaxTokens: 100,
		Temperature: 0.01,
		TopP: 0.01,
	}
	resp, err := client.Query("100 * 20 - 40 = ?", params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Output: " + resp)
	}
}