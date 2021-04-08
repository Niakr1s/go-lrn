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
	flag.Parse()

	var redirectSource redirects.RedirectSource = redirects.DefaultRedirectSource{}
	if *yamlPath != "" {
		redirectSource = redirects.YamlRedirectSource{Path: *yamlPath}
	} else if *jsonPath != "" {
		redirectSource = redirects.JsonRedirectSource{Path: *jsonPath}
	}

	rs, err := redirects.NewRedirects(redirectSource)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server at port 3333")
	http.ListenAndServe(":3333", http.HandlerFunc(redirects.Handler(rs)))
}
