package sources

import (
	"github.com/marmalad3/pokemon/internal/app/sources/pokeapi"
	"github.com/marmalad3/pokemon/internal/app/sources/translation"
)

type ApiSources struct {
	PokeApi        *pokeapi.PokeAPIClient
	TranslationApi *translation.TranslationAPIClient
}
