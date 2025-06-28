package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq"
)

func main() {
    client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
    params := groq.QueryParameters{SystemPrompt: "You are a pirate, answer in pirate style"}
    resp, err := client.Query("Tell me about Go", params)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Response: %s\n", resp)
}