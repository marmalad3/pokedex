package sources

import (
	"github.com/IyadAssaf/poke/internal/app/sources/pokeapi"
	"github.com/IyadAssaf/poke/internal/app/sources/translation"
)

type ApiSources struct {
	PokeApi        *pokeapi.PokeAPIClient
	TranslationApi *translation.TranslationAPIClient
}
