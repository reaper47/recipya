package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	// Dev
	tmpl, err := template.ParseGlob("templates/*gohtml")
	if err != nil {
		fmt.Println(err)
	}
	tmpl.Execute(w, "index.gohtml")

	//config.Tpl.Execute(w, "index.gohtml")
}
