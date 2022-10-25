package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	f, _ := os.Open("../data/ex4.html")
	doc, _ := html.Parse(f)
	var links []Link
	traverse(doc, &links)
	fmt.Println(links)
}

func printNode(node *html.Node) {
	fmt.Printf("Type = %v Data = %q, Attr = %q\n", node.Type, node.Data, node.Attr)
}

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

func traverse(n *html.Node, links *[]Link) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			lnk := Link{Href: getAttr("href", c.Attr), Text: strings.TrimSpace(getText(c))}
			*links = append(*links, lnk)
		} else {
			traverse(c, links)
		}
	}
}
