package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

type MarkdownFormatter struct {
	inCodeBlock bool
	codeStyle   string
	headerStyle string
	boldStyle   string
	resetStyle  string
	quoteStyle  string
	listStyle   string
	language    string
}

func newMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{
		inCodeBlock: false,
		codeStyle:   "\033[36m",   // Cyan
		headerStyle: "\033[1;34m", // Bold blue
		boldStyle:   "\033[1m",    // Bold
		resetStyle:  "\033[0m",    // Reset
		quoteStyle:  "\033[2;37m", // Dim white
		listStyle:   "\033[33m",   // Yellow
		language:    "",
	}
}

const (
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

func (mf *MarkdownFormatter) formatLine(line string) string {
	// Handle code blocks
	if strings.HasPrefix(line, "```") {
		if !mf.inCodeBlock {
			mf.language = strings.TrimPrefix(line, "```")
			mf.inCodeBlock = true
			return fmt.Sprintf("%s%s [%s]%s", mf.codeStyle, line, mf.language, mf.resetStyle)
		}
		mf.inCodeBlock = false
		mf.language = ""
		return fmt.Sprintf("%s%s%s", mf.codeStyle, line, mf.resetStyle)
	}

	if mf.inCodeBlock {
		return fmt.Sprintf("%s%s%s", mf.codeStyle, line, mf.resetStyle)
	}

	// Handle headers with level indication
	if strings.HasPrefix(line, "#") {
		level := len(strings.TrimLeft(line, "#"))
		prefix := strings.Repeat("►", level)
		return fmt.Sprintf("%s%s %s%s", mf.headerStyle, prefix, strings.TrimLeft(line, "# "), mf.resetStyle)
	}

	// Handle numbered lists and bullet points
	if matched, _ := regexp.MatchString(`^\s*[\d]+\.|\s*[\*\-\+]`, line); matched {
		return fmt.Sprintf("%s%s%s", mf.listStyle, line, mf.resetStyle)
	}

	// Handle bold text
	if strings.Contains(line, "**") {
		parts := strings.Split(line, "**")
		for i := 1; i < len(parts); i += 2 {
			if i < len(parts) {
				parts[i] = fmt.Sprintf("%s%s%s", mf.boldStyle, parts[i], mf.resetStyle)
			}
		}
		return strings.Join(parts, "")
	}

	return line
}

func formatMarkdown(text string) string {
	formatter := newMarkdownFormatter()
	lines := strings.Split(text, "\n")
	var formatted []string

	for _, line := range lines {
		formatted = append(formatted, formatter.formatLine(line))
	}

	return strings.Join(formatted, "\n")
}

// Moved the banner to its own function
func printBanner() {
	fmt.Printf("%s%s=== Strike CLI - Ask anything ===%s\n", ColorBold, ColorBlue, ColorReset)
	fmt.Printf("%s► Powered by Groq%s\n", ColorYellow, ColorReset)
	fmt.Printf("%s► Default model: llama-3.3-70b-versatile%s\n", ColorYellow, ColorReset)
	fmt.Printf("Type your question and press Enter (or 'quit' to exit)\n\n")
}

func main() {
	printBanner()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s→%s ", ColorGreen, ColorReset)
		question, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%sError reading input: %v%s\n", ColorRed, err, ColorReset)
			continue
		}

		question = strings.TrimSpace(question)
		if question == "quit" {
			break
		}

		if question == "" {
			continue
		}

		divider := fmt.Sprintf("%s---%s", ColorBlue, ColorReset)
		fmt.Println(divider)

		response, err := makeRequest(question)
		if err != nil {
			fmt.Printf("%sError: %v%s\n", ColorRed, err, ColorReset)
			continue
		}

		fmt.Println(formatMarkdown(response))
		fmt.Println(divider)
	}
}
