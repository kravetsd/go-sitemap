package main

import (
	"encoding/xml"
	"fmt"
	"io"
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

	storage := make(LinkStore)
	store(storage, "https://www.calhoun.io/")
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
		if _, ok := s[l.Href]; !ok {
			if strings.Contains(l.Href, "calhoun.io") && !strings.Contains(l.Href, "mailto") {
				s[l.Href] = l
				store(s, l.Href)
			} else if strings.HasPrefix(l.Href, "/") {
				s[l.Href] = l
				store(s, "https://www.calhoun.io"+l.Href)
			}
		}
	}
}

func getLinks(url string) ([]link.Link, error) {
	body, err := getRespBody(url)
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(body)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func getRespBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
