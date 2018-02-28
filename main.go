package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
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

func saveRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "users.db")
		login := r.FormValue("login")
		password := r.FormValue("password")
		_, err := db.Exec("insert into users (login, password) values (?, ?)", login, password)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		db.Close()
		http.Redirect(w, r, "/", 301)
	} else {
		t, err := template.ParseFiles("templates/registration.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.ExecuteTemplate(w, "registration", nil)
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/registration.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "registration", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "login", nil)
}

func saveLoginHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("login"))
	} else {
		sess.Set("username", r.Form["login"])
		http.Redirect(w, r, "/", 302)
	}
}

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func main() {
	fmt.Println("Listening on port :8000")
	//mux := http.NewServeMux()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/registration", registrationHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/saveregistration", saveRegistrationHandler)
	http.HandleFunc("/savelogin", saveLoginHandler)
	http.HandleFunc("/writelink", writeLinkHandler)
	http.HandleFunc("/addlink", addLink)
	http.HandleFunc("/links", linksHandler)
	http.ListenAndServe(":8000", nil)
}
