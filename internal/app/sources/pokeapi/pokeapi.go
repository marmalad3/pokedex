package pokeapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/IyadAssaf/poke/internal/app/pokedex/models"
	"github.com/IyadAssaf/poke/internal/app/sources/baseapi"
)

type PokeAPIClient struct {
	*baseapi.ApiClient
}

const (
	baseUrlEnvVar  = "POKE_API_URL"
	defaultBaseUrl = "https://pokeapi.co/api/v2"
)

func NewPokeAPIClient(opts ...baseapi.ApiClientOpt) (*PokeAPIClient, error) {
	baseClient, err := baseapi.GenerateNewAPIClient(defaultBaseUrl, baseUrlEnvVar).NewApiClient(opts...)
	if err != nil {
		return nil, err
	}
	return &PokeAPIClient{
		baseClient,
	}, nil
}

// ApiResponse defines how data is returned from the PokeAPI.co REST API
type ApiResponse struct {
	Names        []*ApiResponseNameInLanguage        `json:"names"`
	Descriptions []*ApiResponseDescriptionInLanguage `json:"flavor_text_entries"`
	IsLegendary  bool                                `json:"is_legendary"`
	Habitat      ApiResponseHabitat                  `json:"habitat"`
}

// ApiResponseNameInLanguage returns a pokemon name with a language definition
type ApiResponseNameInLanguage struct {
	Language ApiResponseLanguageDefinition `json:"language"`
	Name     string                        `json:"name"`
}

// ApiResponseDescriptionInLanguage returns a pokemon name with a language definition
type ApiResponseDescriptionInLanguage struct {
	Language    ApiResponseLanguageDefinition `json:"language"`
	Description string                        `json:"flavor_text"`
}

// ApiResponseLanguageDefinition defines a language name
type ApiResponseLanguageDefinition struct {
	Name string `json:"name"`
}

// ApiResponseHabitat defines the pokemon's ApiResponseHabitat
type ApiResponseHabitat struct {
	Name string `json:"name"`
}

// mapToModelForLanguage maps the PokeAPI response to the external models.Pokemon
// based on a language code
func (ar *ApiResponse) mapToModelForLanguage(lang string) *models.Pokemon {
	pokemon := &models.Pokemon{
		IsLegendary: &ar.IsLegendary,
		Habitat:     ar.Habitat.Name,
	}

	for _, nameItem := range ar.Names {
		if strings.EqualFold(nameItem.Language.Name, lang) {
			pokemon.Name = nameItem.Name
			break
		}
	}
	for _, descriptionItem := range ar.Descriptions {
		if strings.EqualFold(descriptionItem.Language.Name, lang) {
			pokemon.Description = descriptionItem.Description
			break
		}
	}

	return pokemon
}

// GetPokemon returns calls the PokeAPI.co server and parses out requested data
// for a particular language
func (p *PokeAPIClient) GetPokemon(ctx context.Context, name string, language string) (*models.Pokemon, error) {
	var responseObj *ApiResponse

	headers := url.Values{}
	headers.Add("Accept", "application/json")

	err := p.DoRequest(ctx, "GET", fmt.Sprintf("pokemon-species/%s", name), nil, headers, &responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj.mapToModelForLanguage(language), nil
}
