package main

import (
	"golang-llms/client"
	"log"
)

func main() {
	modelName := "phi3"

	client, err := client.NewOllamaClient(modelName)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
		return
	}
	client.Run()
}
