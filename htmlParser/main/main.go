package main

import (
	"bytes"
	"flag"
	"fmt"
	"htmlParser"
	"io"
	"log"
	"net/http"
	"regexp"
)

func getFilepath() string {
	filepath := flag.String("file", "../data/ex1.html", "file to be traversed for the links")
	flag.Parse()
	return *filepath
}

var path = regexp.MustCompile("^(http|https)://([a-zA-Z0-9]+).([a-zA-Z0-9.]+)/")

func getDomainName(url string) string {
	m := path.FindStringSubmatch(url)
	if m != nil {
		return m[3]
	} else {
		return ""
	}

}
func main() {
	var allLinks []string
	mapSite("https://www.google.com/", &allLinks)
	fmt.Println(allLinks)
}

func filterSameDomain(domain string, links []htmlParser.Link) []htmlParser.Link {
	var sameDomainLinks []htmlParser.Link
	for _, link := range links {
		dn := getDomainName(link.Href)
		if dn == "" || dn == domain {
			sameDomainLinks = append(sameDomainLinks, link)
		}
	}
	return sameDomainLinks
}

func mapSite(url string, allLinks *[]string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	links := htmlParser.Traverse(bytes.NewReader(b))
	sameDomainLinks := filterSameDomain(getDomainName(url), links)
	for _, sameDomainLink := range sameDomainLinks {
		*allLinks = append(*allLinks, sameDomainLink.Href)
	}
}
