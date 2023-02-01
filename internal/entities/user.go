package entities

import (
	"crypto/rand"
	"encoding/hex"
)

const sizeUserId = 8

func CreateUserId() (string, error) {
	b := make([]byte, sizeUserId)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ValidateUserId(id string) (bool, error) {
	b, err := hex.DecodeString(id)
	if err != nil {
		return false, err
	}
	return len(b) == sizeUserId, nil
}
