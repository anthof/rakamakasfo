package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

func writeLinkHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/writelink.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "writelink", nil)
}

func addLink(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "webdata.db")
		title := r.FormValue("title")
		description := r.FormValue("description")
		url := r.FormValue("url")
		_, err := db.Exec("insert into links (title, description, url) values (?, ?, ?)", title, description, url)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		db.Close()
		http.Redirect(w, r, "/", 301)
	} else {
		t, err := template.ParseFiles("templates/writelink.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.ExecuteTemplate(w, "writelink", nil)
	}
}

func linksHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "webdata.db")

}

func main() {
	fmt.Println("Listening on port :8000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/writelink", writeLinkHandler)
	http.HandleFunc("/addlink", addLink)
	http.HandleFunc("/links", linksHandler)
	http.ListenAndServe(":8000", nil)
}
