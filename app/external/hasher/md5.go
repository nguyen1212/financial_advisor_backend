// Package hasher holds implementation of hash algorithms
package hasher

import (
	"crypto/md5"

	"github.com/financial_advisor/app/services/hasher"
)

type md5Hashing struct{}

func NewMD5() hasher.I {
	return &md5Hashing{}
}

func (h *md5Hashing) Hash(input string) []byte {
	b := md5.Sum([]byte(input))

	return b[:]
}
