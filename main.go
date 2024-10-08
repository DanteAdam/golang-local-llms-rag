package main

import (
	"fmt"
	"golang-llms/llms"
	"golang-llms/rag"
)

func main() {
	// modelName := "phi3"

	// client, err := client.NewOllamaClient(modelName)
	// if err != nil {
	// 	log.Fatalf("Error creating client: %v", err)
	// 	return
	// }
	// client.Run()

	// Initialize the Rag struct with the input PDF file

	rg := rag.NewFileManager("test.pdf")
	text, err := rg.ConvertPdfToText()

	if err != nil {
		fmt.Println("Error:", err)
	}

	store := rag.SaveDocuments(text)
	searchQuery, err := llms.GetUserInput("User: ")

	if err != nil {
		fmt.Println("Error getting user input:", err)
		return
	}

	resDocs, err := rag.Retriever(store, searchQuery)
	if err != nil {
		fmt.Println("Error getting relevant documents:", err)
	}
	
	answer := llms.GetAnswer(resDocs, searchQuery)

	fmt.Println("Answer generated from LLM:", answer)

}
