package cmd

import (
	"html/template"
	"log"
	"net/http"
)

/*func main() {
	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}*/

func index(w http.ResponseWriter, r *http.Request) {
	// Dev
	tmpl, err := template.ParseGlob("templates/*.gohtml")
	template.Must(tmpl.ParseGlob("templates/layout/*.gohtml"))
	template.Must(tmpl.ParseGlob("templates/partials/*.gohtml"))
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, "index.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	// config.Tpl().Execute(w, "index.gohtml")
}
