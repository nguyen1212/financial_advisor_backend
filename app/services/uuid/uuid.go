// Package uuid holds implementations for UUID generation
package uuid

//go:generate mockgen -destination=./mock/mock_$GOFILE -source=$GOFILE -package=mock

type UUIDGenerator interface {
	GetUUID() (string, [16]byte, error)
	Parse(s string) ([16]byte, error)
}
