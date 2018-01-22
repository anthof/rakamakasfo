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
	db, _ := sql.Open("sqlite3", "webdata.db")
	if r.Method == "POST" {
		r.ParseForm()
		title := r.PostFormValue("title")
		description := r.PostFormValue("description")
		url := r.PostFormValue("url")
		tx, _ := db.Begin()
		stmt, _ := tx.Prepare("insert into links (title, description, url) values (?, ?, ?)")
		_, err := stmt.Exec(title, description, url)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		tx.Commit()
		db.Close()
	} else {
		t, err := template.ParseFiles("templates/writelink.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.ExecuteTemplate(w, "writelink", nil)
	}
}

func main() {
	fmt.Println("Listening on port :8000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/writelink", writeLinkHandler)
	http.HandleFunc("/addlink", addLink)
	http.ListenAndServe(":8000", nil)
}
