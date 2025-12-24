package transformer

import (
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

var globalTransformer *mold.Transformer

func init() {
	globalTransformer = modifiers.New()
}

// GetInstance get transformer instance
func GetInstance() *mold.Transformer {
	return globalTransformer
}
