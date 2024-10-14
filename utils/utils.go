package utils

import (
	"errors"
	"fmt"
	"golang-llms/llms"
	"golang-llms/rag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

var userPrompt string

func SetUserPrompt(prompt string) {
	userPrompt = prompt
}

func DeleteFile(fileName string) error {
	filePath := filepath.Join("storage", fileName)

	err := os.Remove(filePath)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func GetFileStorage() (string, error) {
	storagePath := "storage"
	dirEntries, err := os.ReadDir(storagePath)
	if err != nil {
		return "", errors.New("could not access storage folder")
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".pdf") {
			return filepath.Join(storagePath, entry.Name()), nil
		}
	}
	return "", errors.New("no PDF file found in storage")
}
func Setup() *qdrant.Store {
	filePath, err := GetFileStorage()
	if err != nil {
		fmt.Print(err)
	}

	rg := rag.NewFileManager(filePath)

	text, err := rg.ConvertPdfToText()
	if err != nil {
		fmt.Println("Error:", err)
	}
	store := rag.SaveDocuments(text)

	return store
}

func Search(store *qdrant.Store) []schema.Document {
	// searchQuery, err := llms.GetUserInput(userPrompt)
	// if err != nil {
	// 	fmt.Println("Error getting user input:", err)
	// 	return nil
	// }
	searchQuery := userPrompt
	resDocs, err := rag.Retriever(store, searchQuery)
	if err != nil {
		log.Fatal("Error retreving document")
	}
	return resDocs
}

func Execute(resDocs []schema.Document) (string, error) {
	// for {
	// searchQuery, err := llms.GetUserInput(userPrompt)
	// if err != nil {
	// 	fmt.Println("Error getting user input:", err)
	// 	return
	// }
	searchQuery := userPrompt
	// if strings.ToLower(searchQuery) == "bye" {
	// 	fmt.Println("Goodbye!")
	// 	break
	// }
	response, err := llms.GetAnswer(resDocs, searchQuery)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
	// }

}
