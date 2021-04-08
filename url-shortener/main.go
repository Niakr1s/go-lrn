package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Redirects map[string]string

func main() {
	yamlPath := flag.String("yaml", "", "yaml config path")
	flag.Parse()

	var redirectsConfig NewRedirectsConfig = DefaultRedirectsConfig{}
	if *yamlPath != "" {
		redirectsConfig = YamlRedirectsConfig{Path: *yamlPath}
	}

	redirects, err := newRedirects(redirectsConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server at port 3333")
	http.ListenAndServe(":3333", http.HandlerFunc(redirectHandler(redirects)))
}

type NewRedirectsConfig interface{}

type DefaultRedirectsConfig struct{}
type YamlRedirectsConfig struct {
	Path string
}

func newRedirects(newRedirectsConfig NewRedirectsConfig) (Redirects, error) {
	switch t := newRedirectsConfig.(type) {
	case DefaultRedirectsConfig:
		return defaultRedirects(), nil
	case YamlRedirectsConfig:
		return yamlRedirects(t.Path)
	default:
		return nil, fmt.Errorf("not known redirect source")
	}
}

func defaultRedirects() Redirects {
	return map[string]string{
		"google": "http://google.com",
		"ya":     "http://yandex.ru",
		"yandex": "http://yandex.ru",
	}
}

func yamlRedirects(yamlPath string) (Redirects, error) {
	contents, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read yaml file: %v", err)
	}

	redirects := map[string]string{}
	err = yaml.Unmarshal(contents, &redirects)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse yaml file: %v", err)
	}

	return redirects, nil
}

var noRedirectTemplate *template.Template

func newNoRedirectTemplate() *template.Template {
	noRedirectTemplate, err := template.New("noRedirect").Parse(
		`
<body>
<h2>Available redirects</h2>
<ul>
{{ range $shortUrl, $redirectTo := . }}
<li><a href="{{ $shortUrl }}">{{ $shortUrl }}</a></li>
{{ end }}
</ul>
</body>
`)
	if err != nil {
		log.Fatalf("couldn't parse template: %v", err)
	}
	return noRedirectTemplate

}

func init() {
	noRedirectTemplate = newNoRedirectTemplate()
}

func writeNoRedirectTemplate(w http.ResponseWriter, redirects Redirects) {
	err := noRedirectTemplate.Execute(w, redirects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func redirectHandler(redirects Redirects) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("incoming request: ", r.URL.Path)
		shortUrl := strings.TrimLeft(r.URL.Path, "/")
		shortUrl = strings.ToLower(shortUrl)
		if url, ok := redirects[shortUrl]; ok {
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusPermanentRedirect)
		} else {
			writeNoRedirectTemplate(w, redirects)
		}
	}
}
