package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))

	return strings.ReplaceAll(base64.StdEncoding.EncodeToString(h.Sum(nil)), "/", "_")
}
