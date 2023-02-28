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

type LinkStore map[string]link.Link

type Urlset struct {
	URL []Url `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

// func getUrls(ls LinkStore) []Url {
// 	var urls []Url
// 	for _, i := range ls {
// 		urls = append(urls, Url{Loc: i.Href})
// 	}
// 	return urls
// }

func main() {
	// setting the domain va flag
	ln := *flag.String("url", "https://www.calhoun.io/creating-random-strings-in-go/", "url to build sitemap for")

	//requested url
	resp, err := http.Get(ln)

	if err != nil {
		log.Fatal("error requesting url: ", err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	hrefs := filterLinks(hrefs(resp.Body, baseUrl.String()), baseUrl.String())

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

func filterLinks(s []string, base string) []string {
	var hrefs []string
	for _, ln := range s {
		if strings.HasPrefix(ln, base) {
			hrefs = append(hrefs, ln)
		}
	}
	return hrefs
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var hrefs []string
	for _, ln := range links {
		switch {
		case strings.HasPrefix(ln.Href, "/"):
			hrefs = append(hrefs, base+ln.Href)
		case strings.HasPrefix(ln.Href, "http"):
			hrefs = append(hrefs, ln.Href)
		}
	}
	return hrefs
}

// func store(s LinkStore, ln string) {
// 	u := &url.URL{
// 		Scheme: "https",
// 		Host:   ln,
// 	}
// 	fmt.Println(ln, "=>", u.String())
// 	links, err := getLinks(ln)
// 	if err != nil {
// 		fmt.Println("error getting links: ", err)
// 	}
// 	for _, l := range links {
// 		if strings.HasPrefix(l.Href, "#") {
// 			continue
// 		}
// 		if strings.HasPrefix(l.Href, "/") {
// 			// fmt.Println("this is a relative link: ", l.Href)
// 			// fmt.Println("we need to change,", ln, "=>", "https://www.calhoun.io"+l.Href)
// 			l = link.Link{Href: "https://www.calhoun.io" + l.Href, Text: l.Text}
// 		}
// 		//fmt.Println("link: ", l.Href)
// 		if _, ok := s[l.Href]; !ok {
// 			if strings.Contains(l.Href, ln) && !strings.Contains(l.Href, "mailto") {
// 				s[l.Href] = l
// 				store(s, l.Href)
// 			}
// 		}
// 	}
// }

// func getLinks(url string) ([]link.Link, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	links, err := link.Parse(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return links, nil
// }
