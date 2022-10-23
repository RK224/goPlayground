package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"urlshort"
)

func readFile(filePath string) []byte {
	yaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return yaml
}

func main() {
	mux := defaultMux()
	filePath := flag.String("yaml", "../data/redirect.yaml", ".yaml configuration file")
	flag.Parse()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc":      "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":          "https://godoc.org/gopkg.in/yaml.v2",
		"/very-very-important": "https://youtu.be/R8U2ElYYChs",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the jsonHandler using the mapHandler as the
	// fallback
	jsonHandler, err := urlshort.JsonHandler(readFile(*filePath), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
