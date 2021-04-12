package parser_test

import (
	"fmt"
	"lrn/html-link-parser/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHtml(t *testing.T) {
	cases := []struct {
		Html        string
		ErrExpected bool
		Expected    parser.Links
	}{
		{
			`<a href="/dog"></a>`,
			false,
			parser.Links{parser.Link{Href: "/dog", Text: ""}},
		},
		{
			`<a href="/dog">1</a><a href="/dog">2</a>`,
			false,
			parser.Links{
				parser.Link{Href: "/dog", Text: "1"},
				parser.Link{Href: "/dog", Text: "2"},
			},
		},
		{
			`<a href="/dog"><span>Something in a span</span> Text not in a span <b>Bold text!</b></a>`,
			false,
			parser.Links{parser.Link{Href: "/dog", Text: "Something in a span Text not in a span Bold text!"}},
		},
		{
			`<a href="/dog">dog<!-- commented text SHOULD NOT be included! --></a>`,
			false,
			parser.Links{parser.Link{Href: "/dog", Text: "dog"}},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			links, err := parser.ParseHtml(strings.NewReader(c.Html))
			if err != nil {
				assert.True(t, c.ErrExpected)
				return
			}
			assert.Equal(t, c.Expected, links)
		})
	}
}
