package pokedex

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marmalad3/pokemon/internal/app/pokedex/handlers"
	"github.com/marmalad3/pokemon/internal/app/sources"
)

// GetRouter requires source APIs to be provided as dependencies
// returns a HTTP router
func GetRouter(apis *sources.ApiSources) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{name}", handleRequest(apis, handlers.GetOnePokemonHandler))
	r.HandleFunc("/pokemon/translated/{name}", handleRequest(apis, handlers.GetTranslatedPokemonHandler))
	return r
}

// requestHandler defines the signature that handlers should accept.
// The interface{} return parameter will be marshalled as JSON
type requestHandler func(apiSources *sources.ApiSources, r *http.Request) (interface{}, error)

// handleRequest handles returning of HTTP codes based on error presence
// as well as common response headers
func handleRequest(apiSources *sources.ApiSources, handler requestHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respObj, err := handler(apiSources, r)
		if err != nil {
			if err == handlers.ErrStatusNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		respBytes, err := json.Marshal(respObj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}
