## Wordcounter

This project implements the Monte Carlo based unique word counting method described in the paper [*Distinct Elements in Streams: An Algorithm for
the (Text) Book*](https://arxiv.org/abs/2301.10191)

It probabilistically estimates the number of unique words in its input (from either stdin or a filename arg) using a maximum amount of word storage.

### To build/run

This is a [Go language](https://go.dev/) project. Once you have Go installed locally:

To directly run:

`go run main.go <args>`

To build an executable:

`go build -o wordcounter`

### Examples

To run the wordcounter with the default memory size (1000 words) and the input from warandpeace.txt:

`go run main.go warandpeace.txt`

To run the wordcounter with a memory size of 2000 words and the input from warandpeace.txt:

`go run main.go -m 2000 warandpeace.txt`
