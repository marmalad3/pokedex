package handlers

import (
	"fmt"
	"net/http"

	"github.com/IyadAssaf/poke/internal/app/source"
	"github.com/gorilla/mux"
)

const defaultLanguage = "en"

func GetOnePokemonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pokemonName := vars["name"]

	fmt.Println("pokmon name", pokemonName)

	client := source.NewApiClient()
	client.GetPokemon(r.Context(), pokemonName, defaultLanguage)

	w.WriteHeader(http.StatusOK)
}
