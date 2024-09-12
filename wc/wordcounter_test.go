package wc

import (
	"strings"
	"testing"
)

func TestWordReaderRead(t *testing.T) {
	input := "Hello, world! This is a test."
	expectedWords := []string{"hello", "world", "this", "is", "a", "test"}

	reader := NewWordReader(strings.NewReader(input))
	word := make([]byte, 1024)
	var words []string

	for n, err := reader.Read(word); err == nil; n, err = reader.Read(word) {
		words = append(words, string(word[:n]))
	}

	if len(words) != len(expectedWords) {
		t.Fatalf("expected %d words, got %d", len(expectedWords), len(words))
	}

	for i, word := range words {
		if word != expectedWords[i] {
			t.Errorf("expected word %d to be %q, got %q", i, expectedWords[i], word)
		}
	}
}

func TestWordReaderIter(t *testing.T) {
	input := "Hello, world! This is a test."
	expectedWords := []string{"hello", "world", "this", "is", "a", "test"}

	reader := NewWordReader(strings.NewReader(input))
	var words []string

	for word := range reader.Words() {
		words = append(words, word)
	}

	if len(words) != len(expectedWords) {
		t.Fatalf("expected %d words, got %d", len(expectedWords), len(words))
	}

	for i, word := range words {
		if word != expectedWords[i] {
			t.Errorf("expected word %d to be %q, got %q", i, expectedWords[i], word)
		}
	}
}
