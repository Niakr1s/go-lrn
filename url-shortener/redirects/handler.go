package redirects

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

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

func Handler(redirects Redirects) http.HandlerFunc {
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
