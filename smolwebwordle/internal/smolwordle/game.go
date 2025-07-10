package smolwordle

import (
	"fmt"
	"os"
	"strings"
)

type Game struct {
	solution []rune
}

func New(corpus []string) (*Game, error) {
	if len(corpus) == 0 {
		return nil, ErrEmptyCorpus
	}

	return &Game{
		solution: splitToUppercaseCharacters(pickRandomWord(corpus)),
	}, nil
}

const (
	ErrInvalidGuess = gameError("invalid guess length")
)

func (g *Game) Play(guess string) (Feedback, error) {
	err := g.validateGuess(guess)
	if err != nil {
		return Feedback{}, fmt.Errorf("this guess is not the correct length: %w", err)
	}

	characters := splitToUppercaseCharacters(guess)
	fb := computeFeedback(characters, g.solution)
	return fb, nil
}

func (g *Game) validateGuess(guess string) error {
	if len(guess) != len(g.solution) {
		return fmt.Errorf("you guessed a %d word length, remember the answer is %d word length, %w", len(guess), len(g.solution), ErrInvalidGuess)
	}

	return nil
}

func (g *Game) ShowAnswer() string {
	return string(g.solution)
}

func splitToUppercaseCharacters(input string) []rune {
	return []rune(strings.ToUpper(input))
}

func computeFeedback(guess, solution []rune) Feedback {
	result := make(Feedback, len(guess))
	used := make([]bool, len(solution))

	if len(guess) != len(solution) {
		_, _ = fmt.Fprintf(os.Stderr, "guess and solution have different lengths: %d vs %d", len(guess), len(solution))
		return result
	}

	for posInGuess, character := range guess {
		if character == solution[posInGuess] {
			result[posInGuess] = correctPosition
			used[posInGuess] = true
		}
	}

	for posInGuess, character := range guess {
		if result[posInGuess] != absentCharacter {
			continue
		}

		for posInSolution, target := range solution {
			if used[posInSolution] {
				continue
			}

			if character == target {
				result[posInGuess] = wrongPosition
				used[posInSolution] = true
				break
			}
		}
	}

	return result
}
