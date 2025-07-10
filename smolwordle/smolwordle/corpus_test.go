package smolwordle

import (
	"errors"
	"slices"
	"testing"
)

func TestReadCorpus(t *testing.T) {
	tt := map[string]struct {
		file   string
		length int
		err    error
	}{
		"English corpus": {
			file:   "../corpus/english.txt",
			length: 34,
			err:    nil,
		},
		"Empty corpus": {
			file:   "../corpus/empty.txt",
			length: 0,
			err:    ErrCorpusIsEmpty,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			words, err := ReadCorpus(tc.file)
			if !errors.Is(tc.err, err) {
				t.Errorf("expected err %v, got %v", tc.err, err)
			}

			if tc.length != len(words) {
				t.Errorf("expected %d, got %d", tc.length, len(words))
			}
		})
	}
}

func TestPickWord(t *testing.T) {
	corpus := []string{"HELLO", "SALUT", "XAIPE"}
	word := pickWord(corpus)

	if !inCorpus(corpus, word) {
		t.Errorf("expected a word in the corpus, got %q", word)
	}
}

func inCorpus(corpus []string, word string) bool {
	return slices.Contains(corpus, word)
}
