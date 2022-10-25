package main

import (
	"flag"
	"fmt"
	"htmlParser"
	"log"
	"os"
)

func getFilepath() string {
	filepath := flag.String("file", "../data/ex1.html", "file to be traversed for the links")
	flag.Parse()
	return *filepath
}

func main() {
	f, err := os.Open(getFilepath())
	if err != nil {
		log.Fatal(err)
	}
	links := htmlParser.Traverse(f)
	fmt.Println(links)
}
