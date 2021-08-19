package handlers

import (
	"context"
	"net/http"

	"github.com/IyadAssaf/poke/internal/app/pokedex/models"
	"github.com/IyadAssaf/poke/internal/app/sources"
	"github.com/IyadAssaf/poke/internal/app/sources/baseapi"
	"github.com/gorilla/mux"
)

const defaultLanguage = "en"

// GetOnePokemonHandler gets one Pokemon from source APIs and returns
// a *models.Pokemon object to the request handler
func GetOnePokemonHandler(apis *sources.ApiSources, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pokemonName := vars["name"]

	pokemon, err := apis.PokeApi.GetPokemon(r.Context(), pokemonName, defaultLanguage)
	if err != nil {
		apiErr, ok := err.(*baseapi.ApiError)
		if !ok {
			return nil, err
		}
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return nil, ErrStatusNotFound
		default:
			return nil, apiErr
		}
	}

	return pokemon, nil
}

// GetTranslatedPokemonHandler gets one Pokemon from source APIs, applies
// some translation logic defined in translatePokemon and returns a
// a *models.Pokemon object to the request handler
func GetTranslatedPokemonHandler(apis *sources.ApiSources, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pokemonName := vars["name"]

	pokemon, err := apis.PokeApi.GetPokemon(r.Context(), pokemonName, defaultLanguage)
	if err != nil {
		apiErr, ok := err.(*baseapi.ApiError)
		if !ok {
			return nil, err
		}
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			return nil, ErrStatusNotFound
		default:
			return nil, apiErr
		}
	}

	pokemon.Description, err = translatePokemon(r.Context(), apis, pokemon)
	if err != nil {
		return nil, err
	}

	return pokemon, nil
}

func translatePokemon(ctx context.Context, apis *sources.ApiSources, pokemon *models.Pokemon) (string, error) {
	if pokemon.Habitat == "cave" || (pokemon.IsLegendary != nil && *pokemon.IsLegendary) {
		return apis.TranslationApi.TranslateYoda(ctx, pokemon.Description)
	}
	return apis.TranslationApi.TranslateShakespeare(ctx, pokemon.Description)
}
