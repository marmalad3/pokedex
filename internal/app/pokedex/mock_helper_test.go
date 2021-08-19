package pokedex

import (
	"testing"

	"github.com/IyadAssaf/poke/internal/app/pokedex/models"
	"github.com/IyadAssaf/poke/internal/app/sources"
	"github.com/IyadAssaf/poke/internal/app/sources/baseapi"
	"github.com/IyadAssaf/poke/internal/app/sources/pokeapi"
	"github.com/IyadAssaf/poke/internal/app/sources/support"
	"github.com/IyadAssaf/poke/internal/app/sources/translation"
	"github.com/stretchr/testify/assert"
)

func mockAPIs(t *testing.T, pokeApiMock *support.MockInteraction, translateApiMocks []*support.MockInteraction) *sources.ApiSources {
	t.Helper()
	var (
		err  error
		apis = &sources.ApiSources{}
	)

	if pokeApiMock != nil {
		apis.PokeApi, err = pokeapi.NewPokeAPIClient(baseapi.WithHTTPClient(
			support.GetMockHTTPClient(t, pokeApiMock)),
		)
		assert.Nil(t, err)
	}

	if translateApiMocks != nil {
		apis.TranslationApi, err = translation.NewTranslationAPIClient(baseapi.WithHTTPClient(
			support.GetMockHTTPClient(t, translateApiMocks...)),
		)
		assert.Nil(t, err)
	}
	return apis
}

func newPokeAPIResponseObject(t *testing.T, pokemon *models.Pokemon) *pokeapi.ApiResponse {
	var isLegendary bool
	if pokemon.IsLegendary != nil {
		isLegendary = *pokemon.IsLegendary
	}
	return &pokeapi.ApiResponse{
		Names: []*pokeapi.ApiResponseNameInLanguage{
			{
				Language: pokeapi.ApiResponseLanguageDefinition{
					Name: "en",
				},
				Name: pokemon.Name,
			},
		},
		Descriptions: []*pokeapi.ApiResponseDescriptionInLanguage{
			{
				Language: pokeapi.ApiResponseLanguageDefinition{
					Name: "en",
				},
				Description: pokemon.Description,
			},
		},
		IsLegendary: isLegendary,
		Habitat: pokeapi.ApiResponseHabitat{
			Name: pokemon.Habitat,
		},
	}
}

func newTranslatedAPIResponseObject(t *testing.T, translatedText string) *translation.ApiResponse {
	return &translation.ApiResponse{
		Contents: translation.ApiResponseContents{
			Translated: translatedText,
		},
	}
}
