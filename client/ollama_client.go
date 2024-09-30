package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	Client    *api.Client
	ModelName string
}

func (oc *OllamaClient) GetUserInput(prompt string) (string, error) {
	fmt.Printf("%v", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", errors.New("could not get input prompt")
	}

	return strings.TrimSpace(text), nil
}

func (oc *OllamaClient) GenerateResponse(prompt string) error {
	req := &api.GenerateRequest{
		Model:  oc.ModelName, // Use the provided model name
		Prompt: prompt,
		Stream: func() *bool { b := true; return &b }(),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	return oc.Client.Generate(ctx, req, respFunc)
}

func (oc *OllamaClient) Run() error {
	for {
		prompt, err := oc.GetUserInput("User:")

		if err != nil {
			return err
		}

		if strings.ToLower(prompt) == "bye" {
			fmt.Println("Goodbye!")
			break
		}

		err = oc.GenerateResponse(prompt)

		if err != nil {
			panic("failed to generate")
		}
		fmt.Println()
	}
	return nil
}

func NewOllamaClient(modelName string) (*OllamaClient, error) {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		panic("failed to load client")
	}
	return &OllamaClient{
		Client:    client,
		ModelName: modelName,
	}, nil

}
