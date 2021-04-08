package redirects

import (
	"fmt"
	"log"
)

type Redirects map[string]string

type RedirectSource interface{}

type DefaultRedirectSource struct{}
type YamlRedirectSource struct {
	Path string
}

type JsonRedirectSource struct {
	Path string
}

type BoltRedirectSource struct {
	Path string
}

func NewRedirects(redirectSource RedirectSource) (Redirects, error) {
	switch t := redirectSource.(type) {
	case DefaultRedirectSource:
		log.Printf("loading default redirects")
		return defaultRedirects(), nil
	case YamlRedirectSource:
		log.Printf("loading yaml redirects from path %s\n", t.Path)
		return yamlRedirects(t.Path)
	case JsonRedirectSource:
		log.Printf("loading json redirects from path %s\n", t.Path)
		return jsonRedirects(t.Path)
	case BoltRedirectSource:
		log.Printf("loading bolt redirects from path %s\n", t.Path)
		return boltRedirects(t.Path)
	default:
		return nil, fmt.Errorf("not known redirect source")
	}
}
