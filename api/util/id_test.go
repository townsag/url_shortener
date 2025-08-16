package util

import (
	"testing"
	"strings"
)

func TestRandomBase62HappyPath(t *testing.T) {
	length := 10
	result, err := RandomBase62(length)

	if err != nil {
		t.Errorf("received error when creating id: %v", err)
	}
	if len(result) != length {
		t.Errorf("RandomBase62(length) = %s, length of result = %d, want %d", result, len(result), length)
	}
	for _, c := range result {
		if !strings.ContainsRune(base62Chars, c) {
			t.Errorf("RandomBase62(length) = %s, rune %c not in allowed runes: %s", result, c, base62Chars)
		}
	}
}