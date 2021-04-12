package parser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type Links []Link

func ParseHtml(r io.Reader) (Links, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	linkNodes := getAllLinkNodes(doc)
	links := make(Links, len(linkNodes))
	for i, lNode := range linkNodes {
		links[i] = nodeToLink(lNode)
	}

	return links, nil
}

func forAllNodes(root *html.Node, fn func(*html.Node)) {
	node := root.FirstChild
	for node != nil {
		fn(node)
		forAllNodes(node, fn)
		node = node.NextSibling
	}
}

func getAllLinkNodes(node *html.Node) []*html.Node {
	linkNodes := []*html.Node{}

	forAllNodes(node, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			linkNodes = append(linkNodes, n)
		}
	})

	return linkNodes
}

func nodeToLink(n *html.Node) Link {
	l := Link{}

	href, err := getAttribute(n.Attr, "href")
	if err != nil {
		href = ""
	}

	l.Href = href
	l.Text = getInnerText(n)

	return l
}

func getInnerText(node *html.Node) string {
	buf := strings.Builder{}

	forAllNodes(node, func(n *html.Node) {
		if n.Type == html.TextNode {
			fmt.Fprint(&buf, n.Data)
		}
	})

	return buf.String()
}

func getAttribute(attr []html.Attribute, key string) (string, error) {
	for _, a := range attr {
		if a.Key == key {
			return a.Val, nil
		}
	}
	return "", fmt.Errorf("no such attribute")
}
