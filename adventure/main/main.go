package main

import (
	"adventure"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("../tmpl/page.html"))

func handler(w http.ResponseWriter, r *http.Request, story adventure.StoryArc) {
	err := templates.ExecuteTemplate(w, "page.html", story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, adventure.StoryArc)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/view/"):]
		fn(w, r, storyMap[key])
	}
}

var storyMap = adventure.ReadStory("../data/gopher.json")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/intro", http.StatusFound)
	})
	http.HandleFunc("/view/", makeHandler(handler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
