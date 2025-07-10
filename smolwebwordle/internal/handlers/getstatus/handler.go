package getstatus

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/repository"
	"smol/smolwebwordle/internal/session"
)

type gameFinder interface {
	Find(id session.GameID) (session.Game, error)
}

func Handler(finder gameFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue(api.GameID)
		if id == "" {
			http.Error(w, "missing the id of the game", http.StatusBadRequest)
			return
		}

		game, err := finder.Find(session.GameID(id))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, "this game does not exist yet", http.StatusNotFound)
				return
			}

			log.Printf("cannot fetch game %s: %s", id, err)
			http.Error(w, "failed to fetch game", http.StatusInternalServerError)
			return
		}

		apiGame := api.ToGameResponse(game)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(apiGame)
		if err != nil {
			log.Printf("failed to write response: %s", err)
		}
	}
}
