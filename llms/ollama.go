package llms

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/tmc/langchaingo/schema"
)

func GetAnswer(docRetrieved []schema.Document, prompt string) (string,error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic("failed to load client")
	}
	ctx := context.Background()
	// history := memory.NewChatMessageHistory()

	context := ""
	for _, doc := range docRetrieved {
		context += doc.PageContent + "\n"
	}

	promptTemplate := fmt.Sprintf(`
	You are a helpful document assistant. Use the provided context to accurately answer the question below:
	**Context**: %v
	**Question**: %v
	### Constraint:
	- You must restrict your interaction within the confines of provided facts and context.
	- Only use the information provided in the document/context.
	- Do not speculate or add information that is not included in the provided context.
	- Avoid providing personal opinions, or making inferences that aren't based on the provided facts and context.
	- Response in the language compatible with the question.`, context, prompt)

	req := &api.GenerateRequest{
		Model:  "llama3.1", // Use the provided model name
		Prompt: promptTemplate,
		Stream: new(bool),
	}
	var result string
	respFunc := func(resp api.GenerateResponse) error {
		result += resp.Response 
		return nil
	}
	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal("Could not generate the response")
	}
	return result,nil

	// conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

	// fmt.Println(conversation)

	// executor := agents.NewExecutor(
	// 	agents.NewConversationalAgent(llm, nil),
	// agents.WithMemory(conversation),
	// )

	// options := []chains.ChainCallOption{
	// 	chains.WithTemperature(0.8),
	// }

	// res, err := chains.Run(ctx, executor, promptTemplate, options...)
	// if err != nil {
	// 	fmt.Println("Error during chain execution:", err)
	// }
	// return res

}

// func GetUserInput(prompt string) (string, error) {
// 	fmt.Printf(prompt)
// 	reader := bufio.NewReader(os.Stdin)
// 	text, err := reader.ReadString('\n')

// 	if err != nil {
// 		return "", errors.New("could not get input prompt")
// 	}

// 	return strings.TrimSpace(text), nil
// }
