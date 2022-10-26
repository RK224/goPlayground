package main

import (
	"bytes"
	"flag"
	"fmt"
	"htmlParser"
	"io"
	"net/http"
	"regexp"
)

func getFilepath() string {
	filepath := flag.String("file", "../data/ex1.html", "file to be traversed for the links")
	flag.Parse()
	return *filepath
}

func main() {
	siteMap := make(map[string]bool)

	MapSite("https://www.google.com/", &siteMap, 0)
	for key := range siteMap {
		fmt.Println(key)
	}
}

var path = regexp.MustCompile("^(http|https)://([a-zA-Z0-9]+).([a-zA-Z0-9]+).([a-zA-Z0-9.]+)")

func getDomainName(url string) map[string]string {
	m := path.FindStringSubmatch(url)
	if m != nil {
		return map[string]string{"protocol": m[1], "subDomain": m[2], "domain": m[3], "topLevelDomain": m[4]}
	}
	return nil
}
func isSameDomain(targetUrlMap map[string]string, urlMap map[string]string) bool {
	return targetUrlMap["domain"] == urlMap["domain"] && targetUrlMap["topLevelDomain"] == urlMap["topLevelDomain"]
}

func constructUrl(targetUrlMap map[string]string, href string) string {
	return fmt.Sprintf("%v://%v.%v.%v%v", targetUrlMap["protocol"], targetUrlMap["subDomain"], targetUrlMap["domain"], targetUrlMap["topLevelDomain"], href)
}

func filterSameDomain(targetUrlMap map[string]string, links []htmlParser.Link) []htmlParser.Link {
	var sameDomainLinks []htmlParser.Link
	for _, link := range links {
		urlMap := getDomainName(link.Href)
		if urlMap == nil {
			sameDomainLinks = append(sameDomainLinks, htmlParser.Link{Href: constructUrl(targetUrlMap, link.Href), Text: link.Text})
		} else if isSameDomain(targetUrlMap, urlMap) {
			sameDomainLinks = append(sameDomainLinks, link)
		}
	}
	return sameDomainLinks
}

func getResponseReader(url string) io.Reader {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return bytes.NewReader(b)
}

func MapSite(url string, siteMap *map[string]bool, depth int) {
	(*siteMap)[url] = true
	if depth >= 1 {
		return
	}

	links := htmlParser.Traverse(getResponseReader(url))
	sameDomainLinks := filterSameDomain(getDomainName(url), links)
	for _, sameDomainLink := range sameDomainLinks {
		if _, ok := (*siteMap)[sameDomainLink.Href]; ok {
			continue
		}
		MapSite(sameDomainLink.Href, siteMap, depth+1)
	}
}
