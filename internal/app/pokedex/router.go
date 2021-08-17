package pokedex

import (
	"github.com/IyadAssaf/poke/internal/app/pokedex/handlers"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{name}", handlers.GetOnePokemonHandler)
	r.HandleFunc("/pokemon/translated/{name}", handlers.GetOnePokemonHandler)
	return r
}
