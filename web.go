package main

import (
	"net/http"
	"html/template"
)

func main() {
	http.HandleFunc("/", listLogs)
	http.ListenAndServe(":8080", nil)
}

func listLogs(response http.ResponseWriter, r *http.Request) {

	pods := GetPods()

	var templates = template.Must(template.ParseFiles("main.html"))
	templates.ExecuteTemplate(response, "main", &pods)
}




