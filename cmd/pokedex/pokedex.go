package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marmalad3/pokemon/internal/app/pokedex"
	"github.com/marmalad3/pokemon/internal/app/sources"
	"github.com/marmalad3/pokemon/internal/app/sources/pokeapi"
	"github.com/marmalad3/pokemon/internal/app/sources/translation"
)

const defaultPort = "5000"

func main() {
	apis, err := initApis()
	if err != nil {
		// log.Fatalf will exit with code 1
		log.Fatalf("failed to initalise apis: %s", err)
	}

	port := defaultPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	listenAddress := fmt.Sprintf(":%s", port)
	router := pokedex.GetRouter(apis)
	srv := &http.Server{
		Handler: router,
		Addr:    listenAddress,
	}

	log.Printf("listening on %s\n", listenAddress)
	log.Fatal(srv.ListenAndServe())
}

func initApis() (*sources.ApiSources, error) {
	var (
		err  error
		apis = &sources.ApiSources{}
	)
	if apis.PokeApi, err = pokeapi.NewPokeAPIClient(); err != nil {
		return nil, err
	}
	if apis.TranslationApi, err = translation.NewTranslationAPIClient(); err != nil {
		return nil, err
	}
	return apis, nil
}
