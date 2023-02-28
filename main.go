package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kravetsd/link"
)

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	// setting the domain va flag
	ln := flag.String("url", "https://www.calhoun.io/posts", "url to build sitemap for")
	depth := flag.Int("depth", 1, "depth to traverse links")
	flag.Parse()

	// getting the links from the domain
	hrefs := bfs(*ln, *depth)
	for _, ln := range hrefs {
		fmt.Println(ln)
	}

	// creating the xml
	toXml := urlset{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	for _, ln := range hrefs {
		toXml.Urls = append(toXml.Urls, loc{ln})
	}

	x := xml.NewEncoder(os.Stdout)
	x.Indent("", "  ")

	fmt.Print(xml.Header)
	if err := x.Encode(toXml); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}

func bfs(urlString string, depth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlString: struct{}{},
	}

	for i := 0; i <= depth; i++ {
		q, nq = nq, make(map[string]struct{})
		for l, _ := range q {
			if _, ok := seen[l]; ok {
				continue
			}
			seen[l] = struct{}{}
			for _, ln := range get(l) {
				nq[ln] = struct{}{}
			}
		}
	}
	var ret []string
	for k, _ := range seen {
		ret = append(ret, k)
	}
	return ret

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
