package main

import (
	"fmt"
	"golang-llms/filemanager"
	"log"

	"github.com/tmc/langchaingo/llms/huggingface"
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
	embed, err := huggingface.New()
	if err != nil {
		log.Fatal(err)
	}
	fm := filemanager.NewFileManager("test.pdf")
	text, err := fm.ConvertPdfToText()

	if err != nil {
		fmt.Println("Error:", err)
	} 
	store := filemanager.SaveDocuments(text, embed)
}
