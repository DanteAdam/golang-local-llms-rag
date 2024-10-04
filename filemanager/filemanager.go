package filemanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type FileManager struct {
	InputFilePath string
}

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

func SaveDocuments(docs []schema.Document, embedder *embeddings.EmbedderImpl, SenderJid string) *qdrant.Store {
	url, err := url.Parse(os.Getenv("QDRANT_URL"))

	if err != nil {
		log.Fatalf("Error parsing url")
		return nil
	}

	store, err := qdrant.New(
		qdrant.WithURL(*url),
		qdrant.WithAPIKey(os.Getenv("mCRkDRFYqPwih0sDAqn3TMVfHyyKqVatXYS8-cdNISQJmVuku83Zug")),
		qdrant.WithCollectionName(SenderJid),
		qdrant.WithEmbedder(embedder),
	)

	if err != nil {
		log.Fatalf("Error creating qdrant")
		return nil
	}

	_, err = store.AddDocuments(context.Background(), docs)

	if err != nil {
		log.Fatalf("Error adding documents")
		return nil
	}

	fmt.Println("Create store")

	return &store
}

func NewFileManager(inputPath string) *FileManager {
	return &FileManager{
		InputFilePath: inputPath,
	}
}
