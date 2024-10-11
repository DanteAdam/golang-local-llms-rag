package rag

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type FileManager struct {
	InputFilePath string
}

var (
	collectionName = "go-rag"
	QDRANT_URL     = "http://localhost:6333"
	QDRANT_API     = "P1WigS4N-V6YGxoLs-o_qi4wmmgYQL2ttOH6L_-g3xf88I5vXPGJow"
)

func (fm *FileManager) ConvertPdfToText() ([]schema.Document, error) {
	file, err := os.Open(fm.InputFilePath)
	if err != nil {
		return nil, errors.New("error opening file")
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return nil, errors.New("error getting file info")
	}

	doc := documentloaders.NewPDF(file, fileInfo.Size())
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize = 4000
	split.ChunkOverlap = 200
	docs, err := doc.LoadAndSplit(context.Background(), split)

	if err != nil {
		return nil, errors.New("error loading document")
	}

	return docs, nil

}

func Embedder() *embeddings.EmbedderImpl {
	ollamaEmbedderModel, err := ollama.New(
		ollama.WithModel("nomic-embed-text:latest"),
	)

	if err != nil {
		log.Fatalf("Error initializing Ollama: %v", err)
	}

	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaEmbedderModel)

	if err != nil {
		log.Fatalf("Error creating embedder: %v", err)
	}
	return ollamaEmbedder
}

func SaveDocuments(docs []schema.Document) *qdrant.Store {
	embed := Embedder()
	url, err := url.Parse(QDRANT_URL)

	if err != nil {
		log.Fatalf("Error parsing url")
		return nil
	}

	store, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithAPIKey(QDRANT_API),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(embed),
	)
	if err != nil {
		log.Fatalf("Error creating qdrant: %v", err)
		return nil
	}

	_, err = store.AddDocuments(context.Background(), docs)
	if err != nil {
		log.Fatalf("Error adding documents: %v", err)
		return nil
	}

	fmt.Println("Document Processing ....")

	return &store
}

func Retriever(store *qdrant.Store, prompt string) ([]schema.Document, error) {
	vectorOptions := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.5),
	}

	retriever := vectorstores.ToRetriever(store, 10, vectorOptions...)
	docRetrieved, err := retriever.GetRelevantDocuments(context.Background(), prompt)

	if err != nil {
		return nil, fmt.Errorf("could not find relevant information: %v", err)
	}

	return docRetrieved, nil
}

func NewFileManager(inputPath string) *FileManager {
	return &FileManager{
		InputFilePath: inputPath,
	}
}
