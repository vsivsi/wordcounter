package wc

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"math/rand"
	"os"
	"strings"
	"unicode"
)

// WordReader is the struct that implements io.Reader
type WordReader struct {
	scanner *bufio.Scanner
}

// NewWordReader creates a new WordReader
func NewWordReader(r io.Reader) *WordReader {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	return &WordReader{
		scanner: scanner,
	}
}

// Read reads data from the input stream and returns a single lowercase word at a time
func (wr *WordReader) Read(p []byte) (n int, err error) {
	if !wr.scanner.Scan() {
		if err := wr.scanner.Err(); err != nil {
			return 0, err
		}
		return 0, io.EOF
	}
	word := wr.scanner.Text()
	cleanedWord := removeNonAlphabetic(word)
	if len(cleanedWord) == 0 {
		return wr.Read(p)
	}
	n = copy(p, []byte(cleanedWord))
	return n, nil
}

// All returns an iterator allowing the caller to iterate over the WordReader using for/range.
func (wr *WordReader) All() iter.Seq[string] {
	word := make([]byte, 1024)
	return func(yield func(string) bool) {
		var err error
		var n int
		for {
			n, err = wr.Read(word)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if !yield(string(word[:n])) {
				return
			}
		}
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "error reading word: %v\n", err)
		}
	}
}

// removeNonAlphabetic removes non-alphabetic characters from a word using strings.Map
func removeNonAlphabetic(word string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, word)
}

// ProbabilisticSkipper determines if an item should be retained with probability 1/(1<<n)
type ProbabilisticSkipper struct {
	n       int
	counter uint64
	bitmask uint64
}

// NewProbabilisticSkipper initializes the ProbabilisticSkipper
func NewProbabilisticSkipper(n int) *ProbabilisticSkipper {
	pr := &ProbabilisticSkipper{n: n}
	pr.refreshCounter()
	return pr
}

// refreshCounter refreshes the counter with a new random value
func (pr *ProbabilisticSkipper) refreshCounter() {
	if pr.n == 0 {
		pr.bitmask = ^uint64(0) // All bits set to 1
	} else {
		pr.bitmask = rand.Uint64()
		for i := 0; i < pr.n-1; i++ {
			pr.bitmask &= rand.Uint64()
		}
	}
	pr.counter = 64
}

// ShouldSkip returns true with probability 1/(1<<n)
func (pr *ProbabilisticSkipper) ShouldSkip() bool {
	remove := pr.bitmask&1 == 0
	pr.bitmask >>= 1
	pr.counter--
	if pr.counter == 0 {
		pr.refreshCounter()
	}
	return remove
}

// EstimateUniqueWords estimates the number of unique words using a probabilistic counting method
func EstimateUniqueWords(reader io.Reader, memorySize int) int {
	wordReader := NewWordReader(reader)
	words := make(map[string]struct{}, memorySize)

	rounds := 0
	roundRemover := NewProbabilisticSkipper(1)
	wordSkipper := NewProbabilisticSkipper(0)

	for word := range wordReader.All() {
		if wordSkipper.ShouldSkip() {
			delete(words, word)
		} else {
			words[word] = struct{}{}

			if len(words) >= memorySize {
				rounds++
				wordSkipper = NewProbabilisticSkipper(rounds)
				for w := range words {
					if roundRemover.ShouldSkip() {
						delete(words, w)
					}
				}
			}
		}
	}

	if len(words) == 0 {
		return 0
	}

	invProbability := 1 << rounds
	estimatedUniqueWords := len(words) * invProbability
	return estimatedUniqueWords
}

// EstimateUniqueWordsOld estimates the number of unique words using a probabilistic counting method
// EstimateUniqueWords estimates the number of unique words using a probabilistic counting method
func EstimateUniqueWordsOld(reader io.Reader, memorySize int) int {
	wordReader := NewWordReader(reader)
	words := make(map[string]struct{})
	buf := make([]byte, 100)

	rounds := 0
	roundRemover := NewProbabilisticSkipper(1)
	wordSkipper := NewProbabilisticSkipper(0)

	for {
		n, err := wordReader.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		word := string(buf[:n])
		if wordSkipper.ShouldSkip() {
			delete(words, word)
		} else {
			words[word] = struct{}{}

			if len(words) >= memorySize {
				rounds++
				wordSkipper = NewProbabilisticSkipper(rounds)
				for word := range words {
					if roundRemover.ShouldSkip() {
						delete(words, word)
					}
				}
			}
		}
	}

	if len(words) == 0 {
		return 0
	}

	invProbability := 1 << rounds
	estimatedUniqueWords := len(words) * invProbability
	return estimatedUniqueWords
}
