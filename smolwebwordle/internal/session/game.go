package session

import (
	"smol/smolwebwordle/internal/smolwordle"
)

type Game struct {
	ID           GameID
	SmolWordle   smolwordle.Game
	AttemptsLeft byte
	Guesses      []Guess
	Status       Status
}

type GameID string
type Status string

const (
	StatusPlaying = "Playing"
	StatusWon     = "Won"
	StatusLost    = "Lost"
)

type Guess struct {
	Word     string
	Feedback string
}
