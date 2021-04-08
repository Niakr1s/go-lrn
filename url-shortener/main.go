package main

import (
	"encoding/json"
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
	jsonPath := flag.String("json", "", "json config path")
	flag.Parse()

	var redirectSource RedirectSource = DefaultRedirectSource{}
	if *yamlPath != "" {
		redirectSource = YamlRedirectSource{Path: *yamlPath}
	} else if *jsonPath != "" {
		redirectSource = JsonRedirectSource{Path: *jsonPath}
	}

	redirects, err := newRedirects(redirectSource)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server at port 3333")
	http.ListenAndServe(":3333", http.HandlerFunc(redirectHandler(redirects)))
}

type RedirectSource interface{}

type DefaultRedirectSource struct{}
type YamlRedirectSource struct {
	Path string
}

type JsonRedirectSource struct {
	Path string
}

func newRedirects(redirectSource RedirectSource) (Redirects, error) {
	switch t := redirectSource.(type) {
	case DefaultRedirectSource:
		log.Printf("loaded default redirects")
		return defaultRedirects(), nil
	case YamlRedirectSource:
		log.Printf("loaded yaml redirects from path %s\n", t.Path)
		return yamlRedirects(t.Path)
	case JsonRedirectSource:
		log.Printf("loaded json redirects from path %s\n", t.Path)
		return jsonRedirects(t.Path)
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

func jsonRedirects(jsonPath string) (Redirects, error) {
	contents, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read json file: %v", err)
	}

	redirects := map[string]string{}
	err = json.Unmarshal(contents, &redirects)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse json file: %v", err)
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
