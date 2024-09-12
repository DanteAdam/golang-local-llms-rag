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
		log.Fatal(err)
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
			Model:  "llama3",
			Prompt: prompt,

			// set streaming to false
			Stream: new(bool),
		}

		ctx := context.Background()
		respFunc := func(resp api.GenerateResponse) error {
			fmt.Println("Bot:", resp.Response)
			return nil
		}

		err = client.Generate(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}
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
	text = strings.TrimSuffix(text, "\n")

	return text

}
