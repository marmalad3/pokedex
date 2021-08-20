package pokeapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/marmalad3/pokemon/internal/app/sources/baseapi"
	"github.com/marmalad3/pokemon/internal/app/sources/support"
	"github.com/stretchr/testify/assert"
)

func TestApiClientFetchPokemonSuccess(t *testing.T) {
	inputPokemon := ApiResponse{
		Names: []*ApiResponseNameInLanguage{
			{
				Language: ApiResponseLanguageDefinition{
					Name: "en",
				},
				Name: "Oddish",
			},
		},
		Descriptions: []*ApiResponseDescriptionInLanguage{
			{
				Language: ApiResponseLanguageDefinition{
					Name: "en",
				},
				Description: "Oddish is a house plant who enjoys regular watering and occasional misting",
			},
		},
		IsLegendary: true,
		Habitat: ApiResponseHabitat{
			Name: "Iyad's living room",
		},
	}

	respPayload, err := json.Marshal(inputPokemon)
	assert.Nil(t, err)

	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   respPayload,
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/Oddish",
	})

	pokeClient, err := NewPokeAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	foundPokemon, err := pokeClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, err)

	assert.Equal(t, inputPokemon.Names[0].Name, foundPokemon.Name)
	assert.Equal(t, inputPokemon.Descriptions[0].Description, foundPokemon.Description)
	assert.Equal(t, inputPokemon.Habitat.Name, foundPokemon.Habitat)
	assert.Equal(t, inputPokemon.IsLegendary, *foundPokemon.IsLegendary)
}

func TestApiClientFetchPokemonNotFound(t *testing.T) {

	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   []byte(`Not found`),
		ResponseStatus: http.StatusNotFound,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/Iyad",
	})

	pokeClient, err := NewPokeAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	foundPokemon, err := pokeClient.GetPokemon(context.Background(), "Iyad", "en")
	assert.Nil(t, foundPokemon)

	assert.NotNil(t, err)

	assert.IsType(t, &baseapi.ApiError{}, err)

	castErr, ok := err.(*baseapi.ApiError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, castErr.StatusCode)
	assert.Equal(t, "Status '404' returned from server", err.Error())
	assert.False(t, castErr.IsRetryable())
}

func TestApiClientFetchPokemonServerError(t *testing.T) {
	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   []byte(`Internal server error`),
		ResponseStatus: http.StatusInternalServerError,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/Mewtwo",
	})

	pokeClient, err := NewPokeAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	foundPokemon, err := pokeClient.GetPokemon(context.Background(), "Mewtwo", "en")
	assert.Nil(t, foundPokemon)

	assert.NotNil(t, err)
	assert.IsType(t, &baseapi.ApiError{}, err)

	castErr, ok := err.(*baseapi.ApiError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, castErr.StatusCode)
	assert.Equal(t, "Status '500' returned from server", err.Error())

	assert.True(t, castErr.IsRetryable())
}

func TestApiClientFetchPokemonSuccessLangNotFound(t *testing.T) {
	inputPokemon := ApiResponse{
		Names: []*ApiResponseNameInLanguage{
			{
				Language: ApiResponseLanguageDefinition{
					Name: "fr",
				},
				Name: "Mystherbe",
			},
		},
		Descriptions: []*ApiResponseDescriptionInLanguage{
			{
				Language: ApiResponseLanguageDefinition{
					Name: "fr",
				},
				Description: "Oddish est une plante d'intérieur qui aime un arrosage régulier et une brumisation occasionnelle",
			},
		},
		IsLegendary: true,
		Habitat: ApiResponseHabitat{
			Name: "Iyad's living room",
		},
	}

	respPayload, err := json.Marshal(inputPokemon)
	assert.Nil(t, err)

	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   respPayload,
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/Oddish",
	})

	pokeClient, err := NewPokeAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	foundPokemon, err := pokeClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, err)

	assert.Equal(t, "", foundPokemon.Name)
	assert.Equal(t, "", foundPokemon.Description)
	assert.Equal(t, inputPokemon.Habitat.Name, foundPokemon.Habitat)
	assert.Equal(t, inputPokemon.IsLegendary, *foundPokemon.IsLegendary)
}
