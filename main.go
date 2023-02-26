package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kravetsd/link"
)

type Vertex struct {
	Url       string
	IsVisited bool
	Children  []*Vertex
	Parrent   *Vertex
}

func main() {

	root := Vertex{Url: "https://www.calhoun.io", IsVisited: true, Children: []*Vertex{}}

	rootBody, err := getRespBody(root.Url)
	if err != nil {
		log.Println("error parsing url: ", err)
	}

	links, err := link.Parse(rootBody)
	if err != nil {
		log.Println("error parsing links: ", err)
	}
	for _, l := range links {
		root.Children = append(root.Children, &Vertex{Url: l.Href, IsVisited: false, Children: []*Vertex{}})
	}

	fmt.Println(root)

}

func getRespBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
