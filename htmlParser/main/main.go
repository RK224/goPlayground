package main

import (
	"flag"
	"fmt"
	"htmlParser"
)

func getFilepath() string {
	filepath := flag.String("file", "../data/ex1.html", "file to be traversed for the links")
	flag.Parse()
	return *filepath
}

func main() {
	links := htmlParser.Traverse(getFilepath())
	fmt.Println(links)
}
