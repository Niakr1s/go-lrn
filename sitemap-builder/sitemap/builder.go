package sitemap

import (
	"context"
	"fmt"
	"lrn/html-link-parser/parser"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func Build(ctx context.Context, startUrl string, opts ...option) (Sitemap, error) {
	return newBuilder(opts...).BuildSiteMap(ctx, startUrl)
}

type option func(b *builder)

func WithDepth(depth int) option {
	return func(b *builder) {
		b.depth = depth
	}
}

func WithMaxWorkers(limit int) option {
	return func(b *builder) {
		b.sem = make(chan struct{}, limit)
	}
}

type builder struct {
	mu *sync.RWMutex
	wg *sync.WaitGroup

	mem map[string]*url.URL
	sem chan struct{}

	depth int

	errors []error
}

func newBuilder(opts ...option) *builder {
	b := &builder{
		mu:     &sync.RWMutex{},
		wg:     &sync.WaitGroup{},
		mem:    make(map[string]*url.URL),
		sem:    make(chan struct{}, 10),
		depth:  -1,
		errors: make([]error, 0),
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (g *builder) BuildSiteMap(ctx context.Context, startUrl string) (Sitemap, error) {
	u, err := url.Parse(startUrl)
	if err != nil {
		return Sitemap{}, err
	}

	g.processUrl(ctx, u, g.depth)
	g.wg.Wait()

	ret := Sitemap{Links: linkMapToLinkArr(g.mem)}
	return ret, g.getError()
}

func (g *builder) processUrl(ctx context.Context, u *url.URL, depth int) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if depth == 0 {
			return
		}
		if !strings.HasPrefix(u.Scheme, "http") {
			return
		}

		wasSet := g.setUrl(u)
		if !wasSet {
			return
		}

		select {
		case <-ctx.Done():
			return
		case g.sem <- struct{}{}:
			defer func() { <-g.sem }()
		}

		links, err := getLinksFromPage(u.String())
		if err != nil {
			g.addError(err)
			return
		}
		for _, link := range links {
			if isDone(ctx) {
				return
			}
			subUrl, err := makeUrl(u, link.Href)
			if err != nil {
				continue
			}
			if subUrl.Host != u.Host {
				continue
			}

			if g.hasUrl(subUrl) {
				continue
			}

			g.processUrl(ctx, subUrl, depth-1)
		}
	}()
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func (g *builder) hasUrl(u *url.URL) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.mem[u.String()]
	return ok
}

func (g *builder) setUrl(u *url.URL) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	_, ok := g.mem[u.String()]
	g.mem[u.String()] = u
	return !ok
}

func (g *builder) addError(e error) {
	g.errors = append(g.errors, e)
}

func (g *builder) getError() error {
	var err error
	for _, e := range g.errors {
		if err == nil {
			err = e
			continue
		}
		err = fmt.Errorf("%s, %s", err, e)
	}
	return err
}

func makeUrl(parentUrl *url.URL, rawUrl string) (*url.URL, error) {
	subUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	if subUrl.Host == "" {
		subUrl.Host = parentUrl.Host
	}
	if subUrl.Scheme == "" {
		subUrl.Scheme = parentUrl.Scheme
	}
	return subUrl, nil
}

func linkMapToLinkArr(linkMap map[string]*url.URL) []string {
	ret := []string{}

	for _, link := range linkMap {
		ret = append(ret, link.String())
	}

	return ret
}

func getLinksFromPage(url string) (parser.Links, error) {
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	links, err := parser.ParseHtml(res.Body)
	if err != nil {
		return nil, err
	}
	return links, nil
}
