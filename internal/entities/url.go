package entities

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

type URL struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
}

const sizeUrlId = 3

func CreateUrlId() (string, error) {
	b := make([]byte, sizeUrlId)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ValidateUrlId(id string) (bool, error) {
	b, err := hex.DecodeString(id)
	if err != nil {
		return false, err
	}
	return len(b) == sizeUrlId, nil
}

func ValidateUrl(originalUrl string) bool {
	_, domain, _ := strings.Cut(originalUrl, "//")
	if len(domain) == 0 {
		return false
	}
	name, zone, _ := strings.Cut(domain, ".")
	if len(name) == 0 || len(zone) < 2 {
		return false
	}
	return true
}
