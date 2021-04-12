package sitemap

import (
	"encoding/xml"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Links   []string `xml:"url>loc"`
}

type sitemapXml struct {
	XMLName xml.Name `xml:"urlset"`
	Sitemap
}

func (s Sitemap) ToXml() ([]byte, error) {
	sx := sitemapXml{Sitemap: s}
	bytes, err := xml.Marshal(sx)
	if err != nil {
		return nil, err
	}
	res := []byte(xml.Header)
	res = append(res, bytes...)
	return res, nil
}
