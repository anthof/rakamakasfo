package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db, _    = sql.Open("sqlite3", "cache/rakamakasfo.db")
	createDB = "create table if not exist links (title text, description text, url text)"
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

func addLinkHandler(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title")
	description := r.FormValue("description")
	url := r.FormValue("url")

	db.Exec(createDB)
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into links (title, description, url) values (?, ?, ?)")
	_, err := stmt.Exec(title, description, url)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tx.Commit()
	db.Close()
}

func main() {
	fmt.Println("Listening on port :8000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/writelink", writeLinkHandler)
	http.HandleFunc("/addlink", addLinkHandler)
	http.ListenAndServe(":8000", nil)

}
