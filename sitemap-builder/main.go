package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"lrn/sitemap-builder/sitemap"
	"os"
)

func main() {
	siteUrl := flag.String("url", "", "homepage of site to parse")
	flag.Parse()

	sitemap, err := sitemap.Build(context.Background(), *siteUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err := encodeSitemap(os.Stdout, sitemap); err != nil {
		log.Fatal(err)
	}
}

func encodeSitemap(w io.Writer, sitemap sitemap.Sitemap) error {
	fmt.Fprint(w, xml.Header)
	e := xml.NewEncoder(os.Stdout)
	e.Indent("", "  ")
	err := e.Encode(sitemap)

	if err != nil {
		return err
	}

	return nil
}
