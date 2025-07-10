package handlers

import (
	"net/http"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/handlers/getstatus"
	"smol/smolwebwordle/internal/handlers/guess"
	"smol/smolwebwordle/internal/handlers/newgame"
	"smol/smolwebwordle/internal/repository"
)

func NewRouter(db *repository.GameRepository) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc(http.MethodPost+" "+api.NewGameRoute, newgame.Handler(db))
	r.HandleFunc(http.MethodPost+" "+api.GetStatusRoute, getstatus.Handler(db))
	r.HandleFunc(http.MethodPut+" "+api.GuessRoute, guess.Handler(db))

	return r
}
