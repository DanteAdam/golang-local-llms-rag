package utils

import (
	"fmt"
	"golang-llms/llms"
	"golang-llms/rag"
	"log"
	"strings"
)

func Execute() {
	rg := rag.NewFileManager("test.pdf")

	text, err := rg.ConvertPdfToText()
	if err != nil {
		fmt.Println("Error:", err)
	}

	store := rag.SaveDocuments(text)

	for {
		searchQuery, err := llms.GetUserInput("User")
		if err != nil {
			fmt.Println("Error getting user input:", err)
			return
		}
		if strings.ToLower(searchQuery) == "bye" {
			fmt.Println("Goodbye!")
			break
		}
		resDocs, err := rag.Retriever(store, searchQuery)
		if err != nil {
			log.Fatal("Error retreving document")
		}
		fmt.Print("Assistant: ")
		llms.GetAnswer(resDocs, searchQuery)
		fmt.Println()
	}
}
