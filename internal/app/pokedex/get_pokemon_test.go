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

func TestGetOnePokemonSuccess(t *testing.T) {
	inputPokemon := &models.Pokemon{
		Name:        "oddish",
		Description: "awesome",
		Habitat:     "living room",
		IsLegendary: support.BoolPtr(true),
	}

	pokeApiRespObj := newPokeAPIResponseObject(t, inputPokemon)
	pokeRespBytes, err := json.Marshal(pokeApiRespObj)
	assert.Nil(t, err)

	apis := mockAPIs(t, &support.MockInteraction{
		ResponseData:   pokeRespBytes,
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/oddish",
	}, nil)

	server := httptest.NewServer(GetRouter(apis))
	defer server.Close()

	loc, err := url.Parse(server.URL)
	assert.Nil(t, err)

	transport := httptransport.New(loc.Host, "", nil)
	c := operations.New(transport, strfmt.Default)
	resp, err := c.GetPokemonName(&operations.GetPokemonNameParams{
		Context: context.Background(),
		Name:    "oddish",
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	assert.EqualValues(t, inputPokemon, resp.Payload)
}

func TestGetOnePokemonNotFound(t *testing.T) {

	apis := mockAPIs(t, &support.MockInteraction{
		ResponseData:   []byte(""),
		ResponseStatus: http.StatusNotFound,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/ferris",
	}, nil)

	server := httptest.NewServer(GetRouter(apis))
	defer server.Close()

	loc, err := url.Parse(server.URL)
	assert.Nil(t, err)

	transport := httptransport.New(loc.Host, "", nil)
	c := operations.New(transport, strfmt.Default)
	_, err = c.GetPokemonName(&operations.GetPokemonNameParams{
		Context: context.Background(),
		Name:    "ferris",
	})
	assert.NotNil(t, err)

	assert.IsType(t, &operations.GetPokemonNameNotFound{}, err)
}

func TestGetOnePokemonServerError(t *testing.T) {

	apis := mockAPIs(t, &support.MockInteraction{
		// PokeAPI returning invalid JSON
		ResponseData:   []byte(`{ "text": "I'm great at writing valid json`),
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodGet,
		ExpectedPath:   "/api/v2/pokemon-species/oddish",
	}, nil)

	server := httptest.NewServer(GetRouter(apis))
	defer server.Close()

	loc, err := url.Parse(server.URL)
	assert.Nil(t, err)

	transport := httptransport.New(loc.Host, "", nil)
	c := operations.New(transport, strfmt.Default)
	_, err = c.GetPokemonName(&operations.GetPokemonNameParams{
		Context: context.Background(),
		Name:    "oddish",
	})
	assert.NotNil(t, err)

	assert.IsType(t, &operations.GetPokemonNameInternalServerError{}, err)
}
