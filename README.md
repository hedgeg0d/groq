
# [Groq](https://console.groq.com/home) API Library for Go

### Description

`groq` is a lightweight and simple Go library for interacting with the [Groq API](https://console.groq.com/docs/api-reference). It provides an easy-to-use interface to send chat completion requests, including support for streaming responses and customizable query parameters. The library is designed to be minimalistic, making it perfect for quick integration into Go projects.

### Reasons to use

-   **Simplicity**: Intuitive API with straightforward methods like `Ask`, `Query`, and `AskQueryStream`.
-   **Streaming Support**: Handle real-time responses from Groq API using Go channels.
-   **Customizable Parameters**: Fine-tune requests with parameters like `max_tokens`, `temperature`, and `top_p`.
-   **Lightweight**: Minimal dependencies and clean code, ideal for small projects or learning purposes.
-   **Well-Tested**: Includes unit tests to ensure reliability and correctness.

### Installation

To install the library, use the following command:

```bash
go get github.com/ваш_ник/groq-go@v0.1.0
```

Ensure you have a valid Groq API key. It is recommended to set it in the `GROQ_API_KEY` environment variable.

### Examples

Below are examples demonstrating how to use `groq-go` for different use cases.

#### Basic Request with `Ask`

Sending a single chat completion request is just that simple:

```go
	client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
	// uses llama-3.1-8b-instant by default, as the fastest model
	resp, _ := client.Ask("Hello. Tell me about yourself")
	fmt.Println("Output: " + resp)
```
#### Using parameters 
Use `Query` with argument of type `QueryParameters` to specify parameters for LLM.

```go
func main() {
    client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
    params := groq.QueryParameters{
    	Temperature: 0.7,
     	TopP: 0.9,
      	MaxTokens: 1000,
    	SystemPrompt: "You are a pirate, answer in pirate style",
    }
    resp, err := client.Query("Tell me about Go", params)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Response: %s\n", resp)
}
```

#### Streaming Response with `AskQueryStream`

Receive the response in real-time chunks using a Go channel.

```go
package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq-go"
)

func main() {
    client := &groq.GroqClient{
        ApiKey: os.Getenv("GROQ_API_KEY"),
        Model:  "deepseek-r1-distill-llama-70b",
    }
    chunks, err := client.AskQueryStream("Tell me more about Golang.", groq.QueryParameters{})
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

```
For more examples, please, check the `examples` folder.

### Work is in progress

This library is actively being developed. Planned features include:

-   Support for system messages (`role: "system"`) in chat completions.
-   Enhanced error handling for specific Groq API errors (e.g., rate limits).
-   Additional endpoints, such as listing available models (`/models`).
-   More example use cases and documentation.

