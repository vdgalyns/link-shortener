package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
)

const hashLen = 6

func CreateHash(url string) string {
	data := []byte(url)
	sum := fmt.Sprintf("%x", md5.Sum(data))
	hash := strings.Builder{}
	for i, v := range sum {
		if i < hashLen {
			hash.WriteString(string(v))
			continue
		}
		break
	}
	return hash.String()
}

func ValidateHash(hash string) bool {
	return len(hash) == hashLen
}
