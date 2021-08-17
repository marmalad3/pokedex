package main

import (
	"github.com/IyadAssaf/poke/internal/app/pokedex/handlers"
	"github.com/gorilla/mux"
)

func StartServer() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{name}", handlers.GetOnePokemonHandler)

	return r
}
