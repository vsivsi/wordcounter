package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/vsivsi/wordcounter/wc"
)

func main() {
	var memorySize int
	flag.IntVar(&memorySize, "m", 1000, "memory size (words)")
	flag.Parse()

	var reader io.Reader
	args := flag.Args()

	if len(args) > 0 && args[0] != "-" {
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	estimatedUniqueWords := wc.EstimateUniqueWords(reader, memorySize)
	fmt.Printf("Estimated number of unique words: %d\n", estimatedUniqueWords)
}
