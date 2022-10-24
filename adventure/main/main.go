package main

import (
	"adventure"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("../tmpl/view.html"))

func handler(w http.ResponseWriter, r *http.Request, story adventure.StoryArc) {
	err := templates.ExecuteTemplate(w, "view.html", story)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(fn func(http.ResponseWriter, *http.Request, adventure.StoryArc), storyMap map[string]adventure.StoryArc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/view/"):]
		if v, ok := storyMap[key]; ok {
			fn(w, r, v)
		} else {
			http.NotFound(w, r)
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/intro", http.StatusFound)
}

func main() {
	storyMap := adventure.ReadStory("../data/gopher.json")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", viewHandler(handler, storyMap))
	fmt.Println("Starting http server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
