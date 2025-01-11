package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const GROQ_API_URL = "https://api.groq.com/openai/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func makeRequest(question string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	reqBody := Request{
		Model: "llama-3.3-70b-versatile",
		Messages: []Message{
			{Role: "user", Content: question},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", GROQ_API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return response.Choices[0].Message.Content, nil
}

func main() {
	fmt.Println("\033[34m=== Grok CLI - Ask anything ===\033[0m")
	fmt.Println("Type your question and press Enter (or 'quit' to exit)")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n\033[32m→ \033[0m")
		question, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\033[31mError reading input: %v\033[0m\n", err)
			continue
		}

		question = strings.TrimSpace(question)
		if question == "quit" {
			break
		}

		if question == "" {
			continue
		}

		fmt.Println("\033[34m---\033[0m")
		response, err := makeRequest(question)
		if err != nil {
			fmt.Printf("\033[31mError: %v\033[0m\n", err)
			continue
		}

		fmt.Printf("\033[33m%s\033[0m\n", response)
		fmt.Println("\033[34m---\033[0m")
	}
}
