package llms

import (
	"bufio"
	"context"
	"errors"
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
		ollama.WithModel("llama3.1"),
	)
	if err != nil {
		log.Fatalf("could not get ollama service: %v", err)
	}
	return llm
}

func GetAnswer(ctx context.Context, docRetrieved []schema.Document, prompt string) string {
	llm := llmOllama()
	history := memory.NewChatMessageHistory()

	for _, doc := range docRetrieved {
		history.AddAIMessage(ctx, doc.PageContent)
	}

	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

	executor := agents.NewExecutor(
		
		agents.NewConversationalAgent(llm, nil),
		agents.WithMemory(conversation),
	)

	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}

	res, err := chains.Run(ctx, executor, prompt, options...)
	if err != nil {
		fmt.Println("Error running chains", err)
	}

	return res
}

func GetUserInput(prompt string) (string, error) {
	fmt.Printf("%v", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", errors.New("could not get input prompt")
	}

	return strings.TrimSpace(text), nil
}
