package util

import (
	"crypto/rand"
	"math/big"
)

const base62Chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// TODO: benchmark this implementation so I can get a feeling for
//		 wether or not it is slow?
func RandomBase62(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		// for each position in the id, generate a random number
		// between zero and 62. Use the character at that index in the
		// source array as the character at that position in the id
		temp, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[temp.Int64()]
	}
	return string(result), nil
}