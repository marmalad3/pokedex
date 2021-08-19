package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IyadAssaf/poke/internal/app/pokedex"
)

func main() {
	r := pokedex.GetRouter()

	port := "5000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("127.0.0.1:%s", port),
	}

	log.Fatal(srv.ListenAndServe())
}
