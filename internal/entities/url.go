package entities

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type URL struct {
	Hash        string `json:"hash"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
}

const sizeURLHash = 3

func CreateURLHash(originalURL string) (string, error) {
	h := md5.New()
	h.Write([]byte(originalURL))
	s := h.Sum(nil)
	return hex.EncodeToString(s[:sizeURLHash]), nil
}

func ValidateURLHash(hash string) (bool, error) {
	b, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	return len(b) == sizeURLHash, nil
}

func ValidateURL(originalURL string) bool {
	_, domain, _ := strings.Cut(originalURL, "//")
	if len(domain) == 0 {
		return false
	}
	name, zone, _ := strings.Cut(domain, ".")
	if len(name) == 0 || len(zone) < 2 {
		return false
	}
	return true
}
