// Package uuid holds implementation of uuid
package uuid

import (
	"github.com/financial_advisor/app/services/uuid"
	googleUUID "github.com/google/uuid"
)

type uuidv7 struct{}

func NewUUIDv7() uuid.UUIDGenerator {
	return uuidv7{}
}

func (uuidv7) GetUUID() (string, [16]byte, error) {
	value, err := googleUUID.NewV7()

	return value.String(), value, err
}

func (uuidv7) Parse(input string) ([16]byte, error) {
	return googleUUID.Parse(input)
}
