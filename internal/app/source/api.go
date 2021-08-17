package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/IyadAssaf/poke/internal/app/models"
)

var defaultBaseUrl, _ = url.Parse("https://pokeapi.co/api/v2/")

// ApiClient interfaces with the PokeAPI.co REST API
type ApiClient struct {
	client  *http.Client
	baseUrl *url.URL
}

// ApiClientOpt defines optional ApiClient functional params
type ApiClientOpt func(input *ApiClient)

// NewApiClient creates a new ApiClient object using params
func NewApiClient(opts ...ApiClientOpt) *ApiClient {
	client := &ApiClient{
		client:  http.DefaultClient,
		baseUrl: defaultBaseUrl,
	}

	for _, optFn := range opts {
		optFn(client)
	}
	return client
}

// WithHTTPClient is an functional param for setting a http.Client
func WithHTTPClient(c *http.Client) ApiClientOpt {
	return func(as *ApiClient) {
		as.client = c
	}
}

// WithBaseURL is a functional param for setting a base URL for the poke API
func WithBaseURL(b *url.URL) ApiClientOpt {
	return func(as *ApiClient) {
		as.baseUrl = b
	}
}

// apiResponse defines how data is returned from the PokeAPI.co REST API
type apiResponse struct {
	Names        []*nameInLanguage        `json:"names"`
	Descriptions []*descriptionInLanguage `json:"flavor_text_entries"`
	IsLegendary  bool                     `json:"is_legendary"`
	Habitat      habitat                  `json:"habitat"`
}

type nameInLanguage struct {
	Language languageDefinition `json:"language"`
	Name     string             `json:"name"`
}

type descriptionInLanguage struct {
	Language    languageDefinition `json:"language"`
	Description string             `json:"flavour_text"`
}

type languageDefinition struct {
	Name string `json:"name"`
}

type habitat struct {
	Name string `json:"name"`
}

func (ar *apiResponse) mapToModelForLanguage(lang string) *models.Pokemon {
	pokemon := &models.Pokemon{
		IsLegendary: ar.IsLegendary,
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

type ApiError struct {
	StatusCode int
}

func (ae *ApiError) Error() string {
	return fmt.Sprintf("Status '%d' returned from server", ae.StatusCode)
}

func (ae *ApiError) IsRetryable() bool {
	return ae.StatusCode >= 500
}

func (ac *ApiClient) doRequest(ctx context.Context, path string, responseObj interface{}) error {

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", ac.baseUrl, path), nil)
	if err != nil {
		return err
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &ApiError{
			StatusCode: resp.StatusCode,
		}
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(payload, &responseObj); err != nil {
		return err
	}

	return nil
}

// GetPokemon returns calls the PokeAPI.co server and parses out requested  data
func (ac *ApiClient) GetPokemon(ctx context.Context, name string, language string) (*models.Pokemon, error) {
	var responseObj *apiResponse
	err := ac.doRequest(ctx, fmt.Sprintf("pokemon-species/%s", name), &responseObj)
	if err != nil {
		return nil, err
	}

	return responseObj.mapToModelForLanguage(language), nil

}
