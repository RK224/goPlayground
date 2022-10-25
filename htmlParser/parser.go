package htmlParser

import (
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func getAttr(name string, attributes []html.Attribute) string {
	var value string
	for _, attr := range attributes {
		if attr.Key == "href" {
			value = attr.Val
			break
		}
	}
	return value
}

func getText(n *html.Node) string {
	var text string
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		if data := strings.TrimSpace(ch.Data); ch.Type == html.TextNode && data != "" {
			text += (data + " ")
		} else {
			text += getText(ch)
		}
	}
	return text
}

func Traverse(filepath string) []Link {
	var links []Link
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	var root *html.Node
	root, err = html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	traverseHelper(root, &links)
	return links
}
func traverseHelper(n *html.Node, links *[]Link) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			lnk := Link{Href: getAttr("href", c.Attr), Text: strings.TrimSpace(getText(c))}
			*links = append(*links, lnk)
		} else {
			traverseHelper(c, links)
		}
	}
}
