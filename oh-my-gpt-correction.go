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
	// Validate command-line arguments
	validateArguments()

	// Get API key and query from environment variable and command-line argument
	apiKey := os.Getenv("CHATGPT_API_KEY")
	query := os.Args[1]

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Request chat completion from the API
	result, err := requestChatCompletion(client, query)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		os.Exit(1)
	}

	// Clean up the result by removing unnecessary whitespace characters
	trimmed := cleanResult(result)

	// If the result is "ZOMK", exit with status code 3
	if trimmed == "ZOMK" {
		os.Exit(3)
	}

	// Print the comparison between the original query and the result
	printComparison(query, trimmed)
}

// validateArguments checks if the command-line arguments and environment variable are valid
func validateArguments() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		os.Exit(4)
	}

	if os.Getenv("CHATGPT_API_KEY") == "" {
		os.Exit(2)
	}
}

// requestChatCompletion sends a chat completion request to the API and returns the result
func requestChatCompletion(client *openai.Client, query string) (string, error) {
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
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// cleanResult removes unnecessary whitespace characters from the result
func cleanResult(result string) string {
	return regexp.MustCompile(`[\t\r\n]+`).ReplaceAllString(strings.TrimSpace(result), "\n")
}

// printComparison prints the difference between the original query and the result
func printComparison(query, result string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(query, result, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println(result)
}
