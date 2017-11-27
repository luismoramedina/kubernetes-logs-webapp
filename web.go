package main

import "net/http"

func main() {
	http.HandleFunc("/", listLogs)
	http.ListenAndServe(":8080", nil)
}

func listLogs(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("hello"))
}




