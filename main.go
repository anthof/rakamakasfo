package main

import(
	"fmt"
	"net/http"
	"html/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

func writePostHandler(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/writepost.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "writepost", nil)
}



func main()  {
	fmt.Println("Listening on port :8000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/writepost", writePostHandler)
	http.ListenAndServe(":8000", nil)

}
