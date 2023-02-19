package entities

import (
	"crypto/md5"
	"encoding/hex"
)

type Link struct {
	Hash        string `json:"hash"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
	DeletedAt   string `json:"deleted_at"`
}

const sizeLinkHash = 3

func CreateLinkHash(originalURL string) (string, error) {
	h := md5.New()
	h.Write([]byte(originalURL))
	s := h.Sum(nil)
	return hex.EncodeToString(s[:sizeLinkHash]), nil
}

func ValidateLinkHash(hash string) (bool, error) {
	b, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	return len(b) == sizeLinkHash, nil
}
