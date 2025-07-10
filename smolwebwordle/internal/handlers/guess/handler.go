package guess

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/repository"
	"smol/smolwebwordle/internal/session"
	"smol/smolwebwordle/internal/smolwordle"
)

type gameGuesser interface {
	Find(session.GameID) (session.Game, error)
	Update(game session.Game) error
}

func Handler(guesser gameGuesser) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue(api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}

		r := api.GuessRequest{}
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		game, err := guess(session.GameID(id), r.Guess, guesser)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, smolwordle.ErrInvalidGuess):
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errors.Is(err, session.ErrGameOver):
				http.Error(w, err.Error(), http.StatusForbidden)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		apiGame := api.ToGameResponse(game)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func guess(id session.GameID, guess string, db gameGuesser) (session.Game, error) {
	game, err := db.Find(id)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to find game: %w", err)
	}

	if game.AttemptsLeft == 0 || game.Status == session.StatusWon {
		return session.Game{}, session.ErrGameOver
	}

	feedback, err := game.SmolWordle.Play(guess)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to play move: %w", err)
	}

	log.Printf("guess %v is valid in game %s", guess, id)

	game.Guesses = append(game.Guesses, session.Guess{
		Word:     guess,
		Feedback: feedback.String(),
	})

	game.AttemptsLeft -= 1

	switch {
	case feedback.GameWon():
		game.Status = session.StatusWon
	case game.AttemptsLeft == 0:
		game.Status = session.StatusLost
	default:
		game.Status = session.StatusPlaying
	}

	err = db.Update(game)
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to save play: %w", err)
	}
	return game, nil
}
