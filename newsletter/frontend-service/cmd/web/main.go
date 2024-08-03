package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	styles := http.FileServer(http.Dir("./templates/styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))

	script := http.FileServer(http.Dir("./templates/script/"))
	http.Handle("/script/", http.StripPrefix("/script/", script))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w)
	})

	fmt.Println("Starting front end service on port 2000")
	err := http.ListenAndServe(":2000", nil)
	if err != nil {
		log.Panic(err)
	}
}

// This function is used to render the form.
func render(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./templates/form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
