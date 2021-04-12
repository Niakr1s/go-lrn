package main

import (
	"fmt"
	"log"
	"lrn/html-link-parser/parser"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("provide input html file")
	}
	htmlFilePath := os.Args[1]

	links, err := getLinksFromFile(htmlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		fmt.Println(linkToString(link))
	}
}

func getLinksFromFile(filepath string) (parser.Links, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("chouldn't open file: %s", err)
	}
	defer f.Close()

	links, err := parser.ParseHtml(f)
	if err != nil {
		return nil, fmt.Errorf("chouldn't parse file: %s", err)
	}

	return links, nil
}

func linkToString(l parser.Link) string {
	return fmt.Sprintf("%s (%s)", l.Href, l.Text)
}
