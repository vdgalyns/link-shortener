package entities

import (
	"crypto/rand"
	"encoding/hex"
)

type Link struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	OriginalURL string `json:"original_url"`
}

const sizeLinkID = 4

func CreateLinkID() (string, error) {
	b := make([]byte, sizeLinkID)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ValidateLinkID(id string) (bool, error) {
	b, err := hex.DecodeString(id)
	if err != nil {
		return false, err
	}
	return len(b) == sizeLinkID, nil
}
