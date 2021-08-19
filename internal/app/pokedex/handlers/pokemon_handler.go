package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IyadAssaf/poke/internal/app/pokedex/models"
	"github.com/IyadAssaf/poke/internal/app/sources"
	"github.com/IyadAssaf/poke/internal/app/sources/baseapi"
	"github.com/gorilla/mux"
)

const defaultLanguage = "en"

var ErrStatusNotFound error = fmt.Errorf("not found")

func GetOnePokemonHandler(apis *sources.ApiSources, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pokemonName := vars["name"]

	pokemon, err := apis.PokeApi.GetPokemon(r.Context(), pokemonName, defaultLanguage)
	if err != nil {
		err, ok := err.(*baseapi.ApiError)
		if !ok {
			return nil, err
		}
		switch err.StatusCode {
		case http.StatusNotFound:
			return nil, ErrStatusNotFound
		default:
			return nil, err
		}
	}

	return pokemon, nil
}

func GetTranslatedPokemonHandler(apis *sources.ApiSources, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pokemonName := vars["name"]

	pokemon, err := apis.PokeApi.GetPokemon(r.Context(), pokemonName, defaultLanguage)
	if err != nil {
		err, ok := err.(*baseapi.ApiError)
		if !ok {
			return nil, err
		}
		switch err.StatusCode {
		case http.StatusNotFound:
			return nil, ErrStatusNotFound
		default:
			return nil, err
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
