package main

import (
	"log"
	"lrn/choose-your-adventure/adventure"
	"lrn/choose-your-adventure/server"
	"net/http"
	"os"
)

func main() {
	addr := getEnv("ADDR", ":3333")
	adventurePath := getEnv("ADVENTURES_DIR", "")
	loader := adventure.NewFsAdventureLoader(adventurePath)

	s, err := server.NewServer(server.ServerOptions{AdventureLoader: loader})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server started at %s", addr)
	http.ListenAndServe(addr, s)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
