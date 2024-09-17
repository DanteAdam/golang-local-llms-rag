package llms

import (
	"github.com/ollama/ollama/api"
)

type Agent struct {
	Model api.Client

	VectorStore
}
