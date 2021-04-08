package main

import (
	"flag"
	"log"
	"lrn/url-shortener/redirects"
	"net/http"
)

func main() {
	yamlPath := flag.String("yaml", "", "yaml config path")
	jsonPath := flag.String("json", "", "json config path")
	boltPath := flag.String("bolt", "", "bolt config path")
	flag.Parse()

	var redirectSource redirects.RedirectSource = redirects.DefaultRedirectSource{}
	if *yamlPath != "" {
		redirectSource = redirects.YamlRedirectSource{Path: *yamlPath}
	} else if *jsonPath != "" {
		redirectSource = redirects.JsonRedirectSource{Path: *jsonPath}
	} else if *boltPath != "" {
		redirectSource = redirects.BoltRedirectSource{Path: *boltPath}
	}

	rs, err := redirects.NewRedirects(redirectSource)
	if err != nil {
		log.Fatalf("couldn't load redirects: %v", err)
	}
	log.Printf("Redirects loaded")

	log.Printf("Starting server at port 3333")
	http.ListenAndServe(":3333", http.HandlerFunc(redirects.Handler(rs)))
}
