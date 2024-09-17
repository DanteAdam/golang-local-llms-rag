package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

func main() {

	client, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatalf("failed to load client: %v", err)
		return
	}

	for {
		// Get user input
		prompt := getUserInput("User:")

		// Check if the user says "bye" and break the loop
		if strings.ToLower(prompt) == "bye" {
			fmt.Println("Goodbye!")
			break
		}

		req := &api.GenerateRequest{
			Model:  "phi3",
			Prompt: prompt,
			// Stream: new(bool),
			Stream: func() *bool { b := true; return &b }(),
		}

		ctx := context.Background()

		respFunc := func(resp api.GenerateResponse) error {
			fmt.Print(resp.Response)
			return nil
		}

		err = client.Generate(ctx, req, respFunc)
		if err != nil {
			log.Fatalf("failed to generate: %v", err)
			break
		}

		fmt.Println()
	}
}

func getUserInput(prompt string) string {
	fmt.Printf("%v ", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return ""
	}

	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")

	return text
}
