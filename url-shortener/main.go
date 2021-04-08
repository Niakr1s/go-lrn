package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

func main() {
	log.Printf("Starting server at port 3333")
	http.ListenAndServe(":3333", http.HandlerFunc(redirectHandler))
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

var redirects = map[string]string{
	"google": "http://google.com",
	"ya":     "http://yandex.ru",
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming request: ", r.URL.Path)
	shortUrl := strings.TrimLeft(r.URL.Path, "/")
	shortUrl = strings.ToLower(shortUrl)
	if url, ok := redirects[shortUrl]; ok {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusPermanentRedirect)
	} else {
		err := noRedirectTemplate.Execute(w, redirects)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
