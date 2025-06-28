package main

import (
    "fmt"
    "os"
    "github.com/hedgeg0d/groq"
)


func main() {
	client := groq.GroqClient{ApiKey: os.Getenv("GROQ_API_KEY")}
	resp, _ := client.Ask("Hello. Tell me about yourself")
	fmt.Println("Output: " + resp)
}