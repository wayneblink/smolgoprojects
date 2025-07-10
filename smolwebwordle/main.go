package main

import (
	"net/http"
	"smol/smolwebwordle/internal/handlers"
	"smol/smolwebwordle/internal/repository"
)

func main() {
	db := repository.New()

	err := http.ListenAndServe(":8080", handlers.NewRouter(db))
	if err != nil {
		panic(err)
	}
}
