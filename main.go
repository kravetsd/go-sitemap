package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kravetsd/link"
)

func main() {

	client, err := http.Get("https://www.calhoun.io")
	if err != nil {
		log.Println("error openning url: ", err)
	}

	link, err := link.Parse(client.Body)
	if err != nil {
		log.Println("error parsing url: ", err)
	}

	fmt.Println(link)
}
