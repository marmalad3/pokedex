package pokedex

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/IyadAssaf/poke/internal/app/pokedex/client/operations"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestGetOnePokemon(t *testing.T) {
	server := httptest.NewServer(GetRouter())
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
}

func TestGetOnePokemonNotFound(t *testing.T) {
	server := httptest.NewServer(GetRouter())
	defer server.Close()

	loc, err := url.Parse(server.URL)
	assert.Nil(t, err)

	transport := httptransport.New(loc.Host, "", nil)
	c := operations.New(transport, strfmt.Default)
	resp, err := c.GetPokemonName(&operations.GetPokemonNameParams{
		Context: context.Background(),
		Name:    "notfound",
	})
	assert.Nil(t, resp)
	assert.NotNil(t, err)

	_, ok := err.(*operations.GetPokemonNameNotFound)
	assert.True(t, ok)
}
