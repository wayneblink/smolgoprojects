package newgame

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/session"
	"smol/smolwebwordle/internal/smolwordle"

	"github.com/oklog/ulid"
)

const maxAttempts = 5

type gameAdder interface {
	Add(game session.Game) error
}

var corpora = map[string]string{
	"en": "./../../../corpus/english.txt",
}

func Handler(adder gameAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get(api.Lang)
		corpusPath, ok := corpora[lang]
		if !ok {
			corpusPath = corpora["en"]
		}

		game, err := createGame(adder, corpusPath)
		if err != nil {
			log.Printf("unable to create a new game: %s", err)
			http.Error(w, "failed to create a new game", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		apiGame := api.ToGameResponse(game)
		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}

func createGame(adder gameAdder, corpusPath string) (session.Game, error) {
	corpus, err := smolwordle.ReadCorpus("corpus/english.txt")
	if err != nil {
		return session.Game{}, fmt.Errorf("unable to read corpus %w", err)
	}

	game, err := smolwordle.New(corpus)
	if err != nil {
		return session.Game{},
			fmt.Errorf("failed to create a new gordle game")
	}

	g := session.Game{
		ID:           session.GameID(ulid.Now()),
		SmolWordle:   *game,
		AttemptsLeft: maxAttempts,
		Guesses:      []session.Guess{},
		Status:       session.StatusPlaying,
	}

	err = adder.Add(g)
	if err != nil {
		return session.Game{}, fmt.Errorf("failed to save the new game")
	}
	return g, nil
}
