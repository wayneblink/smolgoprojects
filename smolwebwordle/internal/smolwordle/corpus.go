package smolwordle

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	ErrInaccessibleCorpus = corpusError("corpus can't be opened")
	ErrEmptyCorpus        = corpusError("corpus is empty")
)

func ReadCorpus(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %q for reading (%s): %w", path, err, ErrInaccessibleCorpus)
	}

	words := strings.Fields(string(data))

	if len(words) == 0 {
		return nil, ErrEmptyCorpus
	}

	return words, nil
}

func pickRandomWord(corpus []string) string {
	rand.Seed(time.Now().UTC().UnixMicro())
	index := rand.Intn(len(corpus))

	return corpus[index]
}
