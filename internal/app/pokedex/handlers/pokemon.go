package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IyadAssaf/poke/internal/app/source"
	"github.com/gorilla/mux"
)

const defaultLanguage = "en"

func GetOnePokemonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pokemonName := vars["name"]

	client := source.NewApiClient()
	pokemon, err := client.GetPokemon(r.Context(), pokemonName, defaultLanguage)

	if err != nil {
		castErr, ok := err.(*source.ApiError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch castErr.StatusCode {
		case http.StatusNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	b, err := json.Marshal(pokemon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
