package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	apiKey := os.Getenv("CHATGPT_API_KEY")

	query := os.Args[1]
	if query == "" {
		os.Exit(4)
	}

	if apiKey == "" {
		os.Exit(2)
	}
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: `Please fix typos on the following command line and only return the fixed version as answer.
                              If given command is correct, just reply with the word "ZOMK", nothing else:
                              ` + query,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}

	result := resp.Choices[0].Message.Content
	// No newlines pls
	trimmed := regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(result), "\n")
	if trimmed == "ZOMK" {
		// This case is when it could not suggest anything better
		os.Exit(3)
	}

	// Compare query with result
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(query, trimmed, false)
	// Print both lines
	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println(trimmed)
}
