package source

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiClientDefaultValues(t *testing.T) {
	apiClient := NewApiClient()
	assert.Equal(t, defaultBaseUrl, apiClient.baseUrl)
	assert.Equal(t, http.DefaultClient, apiClient.client)
}

func TestApiClientWithBaseURL(t *testing.T) {
	baseURL, err := url.Parse("http://localhost:1234")
	assert.Nil(t, err)

	apiClient := NewApiClient(WithBaseURL(baseURL))

	assert.Equal(t, baseURL, apiClient.baseUrl)
	assert.Equal(t, http.DefaultClient, apiClient.client)
}

func TestApiClientWithClient(t *testing.T) {
	httpClient := &http.Client{}

	apiClient := NewApiClient(WithHTTPClient(httpClient))

	assert.Equal(t, defaultBaseUrl, apiClient.baseUrl)
	assert.Equal(t, httpClient, apiClient.client)
}

func TestApiClientFetchPokemonSuccess(t *testing.T) {
	inputPokemon := apiResponse{
		Names: []*nameInLanguage{
			{
				Language: languageDefinition{
					Name: "en",
				},
				Name: "Oddish",
			},
		},
		Descriptions: []*descriptionInLanguage{
			{
				Language: languageDefinition{
					Name: "en",
				},
				Description: "Oddish is a house plant who enjoys regular watering and occasional misting",
			},
		},
		IsLegendary: true,
		Habitat: habitat{
			Name: "Iyad's living room",
		},
	}

	respPayload, err := json.Marshal(inputPokemon)
	assert.Nil(t, err)

	httpClient := &http.Client{}
	httpClient.Transport = getMockRoundTripper(t, respPayload, http.StatusOK)

	apiClient := NewApiClient(WithHTTPClient(httpClient))

	foundPokemon, err := apiClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, err)

	assert.Equal(t, inputPokemon.Names[0].Name, foundPokemon.Name)
	assert.Equal(t, inputPokemon.Descriptions[0].Description, foundPokemon.Description)
	assert.Equal(t, inputPokemon.Habitat.Name, foundPokemon.Habitat)
	assert.Equal(t, inputPokemon.IsLegendary, foundPokemon.IsLegendary)
}

func TestApiClientFetchPokemonNotFound(t *testing.T) {
	httpClient := &http.Client{}
	httpClient.Transport = getMockRoundTripper(t, []byte(`Not found`), http.StatusNotFound)

	apiClient := NewApiClient(WithHTTPClient(httpClient))

	foundPokemon, err := apiClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, foundPokemon)

	assert.NotNil(t, err)

	assert.IsType(t, &ApiError{}, err)

	castErr, ok := err.(*ApiError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, castErr.StatusCode)
	assert.Equal(t, "Status '404' returned from server", err.Error())
	assert.False(t, castErr.IsRetryable())
}

func TestApiClientFetchPokemonServerError(t *testing.T) {
	httpClient := &http.Client{}
	httpClient.Transport = getMockRoundTripper(t, []byte(`Internal server error`), http.StatusInternalServerError)

	apiClient := NewApiClient(WithHTTPClient(httpClient))

	foundPokemon, err := apiClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, foundPokemon)

	assert.NotNil(t, err)
	assert.IsType(t, &ApiError{}, err)

	castErr, ok := err.(*ApiError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, castErr.StatusCode)
	assert.Equal(t, "Status '500' returned from server", err.Error())

	assert.True(t, castErr.IsRetryable())
}

func TestApiClientFetchPokemonSuccessLangNotFound(t *testing.T) {
	inputPokemon := apiResponse{
		Names: []*nameInLanguage{
			{
				Language: languageDefinition{
					Name: "fr",
				},
				Name: "Mystherbe",
			},
		},
		Descriptions: []*descriptionInLanguage{
			{
				Language: languageDefinition{
					Name: "fr",
				},
				Description: "Oddish est une plante d'intérieur qui aime un arrosage régulier et une brumisation occasionnelle",
			},
		},
		IsLegendary: true,
		Habitat: habitat{
			Name: "Iyad's living room",
		},
	}

	respPayload, err := json.Marshal(inputPokemon)
	assert.Nil(t, err)

	httpClient := &http.Client{}
	httpClient.Transport = getMockRoundTripper(t, respPayload, http.StatusOK)

	apiClient := NewApiClient(WithHTTPClient(httpClient))

	foundPokemon, err := apiClient.GetPokemon(context.Background(), "Oddish", "en")
	assert.Nil(t, err)

	assert.Equal(t, "", foundPokemon.Name)
	assert.Equal(t, "", foundPokemon.Description)
	assert.Equal(t, inputPokemon.Habitat.Name, foundPokemon.Habitat)
	assert.Equal(t, inputPokemon.IsLegendary, foundPokemon.IsLegendary)
}
