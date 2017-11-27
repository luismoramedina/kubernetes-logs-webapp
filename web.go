package main

import (
	"net/http"
	"html/template"
	"fmt"
	"github.com/jhoonb/archivex"
	"sync"
	"os"
)

func main() {
	http.HandleFunc("/", listLogs)
	http.HandleFunc("/download", download)
	http.ListenAndServe(":8080", nil)
}

func listLogs(response http.ResponseWriter, request *http.Request) {

	pods := GetPods()

	var templates = template.Must(template.ParseFiles("main.html"))
	templates.ExecuteTemplate(response, "main", &pods)
}

func download(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := request.Form
	fmt.Println(form)

	filter := form.Get("filter")
	fmt.Println(filter)

	form.Del("filter")

	keys := make([]string, 0, len(form))
	for k := range form {
		keys = append(keys, k)
	}

	zipFile := getLogsAsync(keys)
	defer os.Remove(zipFile)
	http.ServeFile(response, request, zipFile)

}

func getLogsAsync(pods []string) string {
	zip := new(archivex.ZipFile)
	zip.Create("logs")

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(pods))

	for _, pod := range pods {
		go func(pod string) {
			defer waitGroup.Done()
			fmt.Println(pod)
			logs := getLogs(pod)
			zip.Add(pod + ".txt" , []byte(logs))
		}(pod)

	}
	waitGroup.Wait()
	zip.Close()
	return zip.Name
}




