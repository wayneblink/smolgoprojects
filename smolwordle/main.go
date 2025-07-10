package main

import (
	"bufio"
	"fmt"
	"os"
	"smolwordle/smolwordle"
)

const maxAttempts = 6

func main() {
	corpus, err := smolwordle.ReadCorpus("corpus/english.txt")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to read corpus: %s", err)
	}

	g, err := smolwordle.New(bufio.NewReader(os.Stdin), corpus, maxAttempts)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to start game: %s", err)
		return
	}

	g.Play()
}
