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

func GetAnswer(docRetrieved []schema.Document, prompt string) string {

	llm := llmOllama()
	ctx := context.Background()
	// history := memory.NewChatMessageHistory()

	context := ""
	for _, doc := range docRetrieved {
		context = doc.PageContent + "/n"
	}

	promptTemplate := fmt.Sprintf(`
	You are a helpful document assistant. Use the provided context to accurately answer the question below:
	Context: %v
	Question: %v
	Constraint:
	- Your response must be in Vietnamese.
	- Do no hallucination.`, context, prompt)

	// conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

	// fmt.Println(conversation)

	executor := agents.NewExecutor(
		agents.NewConversationalAgent(llm, nil),
		// agents.WithMemory(conversation),
	)

	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}

	res, err := chains.Run(ctx, executor, promptTemplate, options...)
	if err != nil {
		fmt.Println("Error during chain execution:", err)
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
