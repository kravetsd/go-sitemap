package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kravetsd/link"
)

type LinkStore map[string]link.Link

func main() {

	storage := make(LinkStore)

	store(storage, "https://www.calhoun.io/")
	for _, i := range storage {
		log.Println(i.Href)
	}
	fmt.Println(len(storage))

}

func store(s LinkStore, ln string) {
	body, err := getRespBody(ln)
	if err != nil {
		log.Println("error openning url: ", err)
	}
	links, err := link.Parse(body)
	if err != nil {
		log.Println("error patsing web page: ", err)
	}
	for _, l := range links {
		if _, ok := s[l.Href]; !ok {
			if strings.Contains(l.Href, "calhoun.io") && !strings.Contains(l.Href, "mailto") {
				s[l.Href] = l
				store(s, l.Href)
			} else if strings.HasPrefix(l.Href, "/") {
				s[l.Href] = l
				store(s, "https://www.calhoun.io"+l.Href)
			} else {
				fmt.Println("duplicate or not a calhoun.io link: ", l.Href)
			}
		}
	}
}

func getRespBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
