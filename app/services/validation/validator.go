// Package validation provides a singleton instance of a validator
package validation

import (
	"github.com/go-playground/validator/v10"
)

var globalValidator *validator.Validate

func init() {
	globalValidator = validator.New()
}

// GetInstance global validator
func GetInstance() *validator.Validate {
	return globalValidator
}
