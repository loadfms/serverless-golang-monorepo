package utils

import (
	"crypto/sha256"
	"fmt"
)

func GenerateSHA256(value string) string {
	h := sha256.New()

	h.Write([]byte(value))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
