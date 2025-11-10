package models

import (
	"math/rand"
	"time"
)

const (
	chars  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length = 6
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateShortCode() string {
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
