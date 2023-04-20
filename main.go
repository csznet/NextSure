package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"nextsure/snapshot"
	"regexp"
)

func main() {
	url := "https://www.csz.net"
	snapshot.Get(url)
	println(title(url))
}

func title(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	title := re.FindStringSubmatch(string(body))[1]

	return title
}
