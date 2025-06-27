package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq"
)

func main() {
    client := &groq.GroqClient{
        ApiKey: os.Getenv("GROQ_API_KEY"),
        Model:  "deepseek-r1-distill-llama-70b",
    }

    chunks, err := client.AskQueryStream("Tell me about Golang, i want a lot of information", groq.QueryParameters{})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Println("Streaming response:")
    for chunk := range chunks {
        fmt.Print(chunk)
    }
    fmt.Printf("\nRequests made: %d\n", client.RequestsCount)
}