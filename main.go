package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Link struct {
	Id    int
	Title string
	Url   string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

func writeLinkHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/writelink.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "writelink", nil)
}

func addLink(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "webdata.db")
		title := r.FormValue("title")
		url := r.FormValue("url")
		_, err := db.Exec("insert into links (title, url) values (?, ?)", title, url)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		db.Close()
		http.Redirect(w, r, "/", 301)
	} else {
		t, err := template.ParseFiles("templates/writelink.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.ExecuteTemplate(w, "writelink", nil)
	}
}

func linksHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "webdata.db")
	rows, _ := db.Query("select * from links")
	defer rows.Close()
	links := make([]*Link, 0)

	for rows.Next() {
		link := new(Link)
		err := rows.Scan(&link.Id, &link.Title, &link.Url)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		links = append(links, link)

	}
	t, _ := template.ParseFiles("templates/links.html", "templates/header.html", "templates/footer.html")
	t.ExecuteTemplate(w, "links", links)
}

func main() {
	fmt.Println("Listening on port :8000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/writelink", writeLinkHandler)
	http.HandleFunc("/addlink", addLink)
	http.HandleFunc("/links", linksHandler)
	http.ListenAndServe(":8000", nil)
}
