package pokedex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IyadAssaf/poke/internal/app/pokedex/client/operations"
	"github.com/IyadAssaf/poke/internal/app/pokedex/models"
	"github.com/IyadAssaf/poke/internal/app/sources/support"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestPokemonTranslation(t *testing.T) {

	testCases := []struct {
		name                   string
		pokemon                *models.Pokemon
		yodaTranslation        string
		shakespeareTranslation string
		expectedDescription    string
	}{
		{
			name: "Not legendary and habitat isn't cave, should have shakespeare translation",
			pokemon: &models.Pokemon{
				Name:        "oddish",
				Description: "oddish is a plant",
				Habitat:     "living room",
				IsLegendary: support.BoolPtr(false),
			},
			shakespeareTranslation: "oddish isth a planteth",
			expectedDescription:    "oddish isth a planteth",
		},
		{
			name: "Not legendary but habitat is cave, should have yoda translation",
			pokemon: &models.Pokemon{
				Name:        "oddish",
				Description: "oddish is a plant",
				Habitat:     "cave",
				IsLegendary: support.BoolPtr(false),
			},
			yodaTranslation:     "a plant, oddish is",
			expectedDescription: "a plant, oddish is",
		},
		{
			name: "Is legendary and habitat isn't cave, should have yoda translation",
			pokemon: &models.Pokemon{
				Name:        "oddish",
				Description: "oddish is a plant",
				Habitat:     "living room",
				IsLegendary: support.BoolPtr(true),
			},
			yodaTranslation:     "a plant, oddish is",
			expectedDescription: "a plant, oddish is",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pokeRespBytes, err := json.Marshal(newPokeAPIResponseObject(t, tc.pokemon))
			assert.Nil(t, err)

			pokeMockData := &support.MockInteraction{
				ResponseData:   pokeRespBytes,
				ResponseStatus: http.StatusOK,
				ExpectedMethod: http.MethodGet,
				ExpectedPath:   "/api/v2/pokemon-species/oddish",
			}

			yodaRespBytes, err := json.Marshal(newTranslatedAPIResponseObject(t, tc.yodaTranslation))
			assert.Nil(t, err)
			yodaMockData := &support.MockInteraction{
				ResponseData:   yodaRespBytes,
				ResponseStatus: http.StatusOK,
				ExpectedMethod: http.MethodPost,
				ExpectedPath:   "/translate/yoda.json",
			}

			shakespeareRespBytes, err := json.Marshal(newTranslatedAPIResponseObject(t, tc.shakespeareTranslation))
			assert.Nil(t, err)
			shakespeareMockData := &support.MockInteraction{
				ResponseData:   shakespeareRespBytes,
				ResponseStatus: http.StatusOK,
				ExpectedMethod: http.MethodPost,
				ExpectedPath:   "/translate/shakespeare.json",
			}
			apis := mockAPIs(t, pokeMockData, []*support.MockInteraction{yodaMockData, shakespeareMockData})

			server := httptest.NewServer(GetRouter(apis))
			defer server.Close()

			loc, err := url.Parse(server.URL)
			assert.Nil(t, err)

			transport := httptransport.New(loc.Host, "", nil)
			c := operations.New(transport, strfmt.Default)
			resp, err := c.GetPokemonTranslatedName(&operations.GetPokemonTranslatedNameParams{
				Context: context.Background(),
				Name:    tc.pokemon.Name,
			})
			assert.Nil(t, err)
			assert.NotNil(t, resp)

			assert.Equal(t, tc.expectedDescription, resp.Payload.Description)
		})
	}
}
