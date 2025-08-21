package utils

import (
	"golang.org/x/crypto/sha3"
)

func Keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}