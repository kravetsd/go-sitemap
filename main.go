package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kravetsd/link"
)

type LinkStore map[string]link.Link

type Urlset struct {
	URL []Url `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

func getUrls(ls LinkStore) []Url {
	var urls []Url
	for _, i := range ls {
		urls = append(urls, Url{Loc: i.Href})
	}
	return urls
}

func main() {
	ln := *flag.String("url", "https://www.calhoun.io", "url to build sitemap for")

	storage := make(LinkStore)
	store(storage, ln)
	urlset := Urlset{URL: getUrls(storage)}

	res, _ := xml.MarshalIndent(urlset, "", "  ")
	res = append([]byte(xml.Header), res...)

	f, err := os.Create("sitemap.xml")
	if err != nil {
		log.Fatal("error creating xml file", err)
	}
	f.Write(res)

}

func store(s LinkStore, ln string) {
	links, err := getLinks(ln)
	if err != nil {
		fmt.Println("error getting links: ", err)
	}
	for _, l := range links {
		if strings.HasPrefix(l.Href, "#") {
			continue
		}
		if strings.HasPrefix(l.Href, "/") {
			l = link.Link{Href: "https://www.calhoun.io" + l.Href, Text: l.Text}
		}
		if _, ok := s[l.Href]; !ok {
			if strings.Contains(l.Href, "calhoun.io") && !strings.Contains(l.Href, "mailto") {
				s[l.Href] = l
				store(s, l.Href)
			}
		}
	}
}

func getLinks(url string) ([]link.Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	links, err := link.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return links, nil
}
