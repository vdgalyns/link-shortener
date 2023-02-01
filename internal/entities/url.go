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

const sizeUrlHash = 3

func CreateUrlHash(originalUrl string) (string, error) {
	h := md5.New()
	h.Write([]byte(originalUrl))
	s := h.Sum(nil)
	return hex.EncodeToString(s[:sizeUrlHash]), nil
}

func ValidateUrlHash(hash string) (bool, error) {
	b, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	return len(b) == sizeUrlHash, nil
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
