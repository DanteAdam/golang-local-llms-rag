package llms

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
)

func llmOllama() *ollama.LLM {
	llm, err := ollama.New(
		ollama.WithModel("llama3"),
	)
	if err != nil {
		log.Fatalf("could not get ollama service: %v", err)
	}
	return llm
}

func GetAnswer(ctx context.Context, docRetrieved []schema.Document, prompt string) (string, error) {
	llm := llmOllama()
	history := memory.NewChatMessageHistory()

	for _, doc := range docRetrieved {
		history.AddAIMessage(ctx, doc.PageContent)
	}

	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))
	executor := agents.NewExecutor(
		agents.NewConversationalAgent(llm, nil),
		nil,
		agents.WithMemory(conversation),
	)

	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}

	res, err := chains.Run(ctx, executor, prompt, options...)
	if err != nil {
		return "", err
	}
	return res, nil
}

func GetUserInput(promptString string) (string, error) {
	fmt.Print(promptString, ": ")

	reader := bufio.NewReader(os.Stdin)

	Input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("Could not get the user prompt")
	}

	Input = strings.TrimSuffix(Input, "\n")
	Input = strings.TrimSuffix(Input, "\r")

	return Input, nil
}
