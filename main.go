package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kravetsd/link"
)

func main() {
	// setting the domain va flag
	ln := flag.String("url", "https://www.calhoun.io/creating-random-strings-in-go/", "url to build sitemap for")
	flag.Parse()

	fmt.Println("Getting links from: ", *ln)
	// getting the links from the domain
	hrefs := get(*ln)
	for _, ln := range hrefs {
		fmt.Println(ln)
	}

	// storage := make(LinkStore)
	// store(storage, ln)
	// urlset := Urlset{URL: getUrls(storage)}

	// res, _ := xml.MarshalIndent(urlset, "", "  ")
	// res = append([]byte(xml.Header), res...)

	// Creating an xml file
	// f, err := os.Create("sitemap.xml")
	// if err != nil {
	// 	log.Fatal("error creating xml file", err)
	// }
	// f.Write(res)

}

func filterLinks(s []string, keepFunc func(s string) bool) []string {
	var hrefs []string
	for _, ln := range s {
		if keepFunc(ln) {
			hrefs = append(hrefs, ln)
		}
	}
	return hrefs
}

func withPrefix(pfx string) func(s string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, pfx)
	}
}

func get(urlString string) []string {
	//requested url
	resp, err := http.Get(urlString)

	if err != nil {
		log.Fatal("error requesting url: ", err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	return filterLinks(hrefs(resp.Body, baseUrl.String()), withPrefix(baseUrl.String()))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, ln := range links {
		switch {
		case strings.HasPrefix(ln.Href, "/"):
			ret = append(ret, base+ln.Href)
		case strings.HasPrefix(ln.Href, "http"):
			ret = append(ret, ln.Href)
		}
	}
	return ret
}
