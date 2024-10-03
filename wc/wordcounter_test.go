package wc

import (
	"strings"
	"testing"
)

var (
	input         = "Hello, world! This is a test."
	expectedWords = []string{"hello", "world", "this", "is", "a", "test"}
)

func TestWordReaderRead(t *testing.T) {

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

	reader := NewWordReader(strings.NewReader(input))
	var words []string

	for word := range reader.All() {
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

func TestEstimateUniqueWords(t *testing.T) {
	expectedUniqueWords := 6

	reader := strings.NewReader(input)
	memorySize := 4

	estimatedUniqueWords := EstimateUniqueWords(reader, memorySize)

	if estimatedUniqueWords != expectedUniqueWords {
		t.Errorf("expected %d unique words, got %d", expectedUniqueWords, estimatedUniqueWords)
	}
}

func TestEstimateUniqueWordsOld(t *testing.T) {
	expectedUniqueWords := 6 // "hello", "world", "this", "is", "a", "test" (but "hello" and "world" are repeated)

	reader := strings.NewReader(input)
	memorySize := 6

	estimatedUniqueWords := EstimateUniqueWordsOld(reader, memorySize)

	if estimatedUniqueWords != expectedUniqueWords {
		t.Errorf("expected %d unique words, got %d", expectedUniqueWords, estimatedUniqueWords)
	}
}
