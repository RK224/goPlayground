package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"urlshort"
)

var validPath = regexp.MustCompile(".(yaml|json)$")

func readConf(filePath string) ([]byte, string) {
	m := validPath.FindStringSubmatch(filePath)
	if m == nil {
		log.Fatal("Unsupported filetype : " + filePath)
	}
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return conf, m[1]
}

func readUrlShortParameters() (filePath *string) {
	filePath = flag.String("conf", "../data/redirect.yaml", "configuration file to be used, .json and .yaml files are supported")
	flag.Parse()
	return
}
func main() {
	mux := defaultMux()

	filePath := readUrlShortParameters()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc":      "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":          "https://godoc.org/gopkg.in/yaml.v2",
		"/very-very-important": "https://youtu.be/R8U2ElYYChs",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	conf, ext := readConf(*filePath)
	handler, err := urlshort.Handler(conf, ext, mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
