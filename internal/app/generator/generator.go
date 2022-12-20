package generator

import "math/rand"

func Make() string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := 0; i < len(b); i++ {
		b[i] = charset[rand.Intn(len(charset)-1)]
	}
	return string(b)
}
