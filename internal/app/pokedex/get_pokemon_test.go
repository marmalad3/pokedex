package main

import (
	"context"
	"os"
	"testing"

	"github.com/IyadAssaf/poke/internal/app/client/operations"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestGetOnePokemon(t *testing.T) {
	transport := httptransport.New(os.Getenv("API_HOST"), "", nil)

	// create the API client, with the transport
	c := operations.New(transport, strfmt.Default)

	resp, err := c.GetPokemonName(&operations.GetPokemonNameParams{
		Context: context.Background(),
		Name:    "Oddish",
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
