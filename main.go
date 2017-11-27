package main

import (
   "io/ioutil"
   "github.com/go-errors/errors"
   "github.com/jhoonb/archivex"
   "github.com/Jeffail/gabs"
   "crypto/tls"
   "fmt"
   "strings"
   "net/http"
   "net/url"
)

var namespace = "isb-npccd-dev"

func init() {
   //TODO create client
}

func main() {
   zip := new(archivex.ZipFile)
   zip.Create("logs")
   x := getPods()
   //fmt.Println(x)
   for _, pod := range x {
      fmt.Println(pod)
      logs := getLogs(pod)
      zip.Add(pod + ".txt" , []byte(logs))
   }
   zip.Close()

}

func getLogs(podName string) string {
   client := getClient()

   myurl := getUrl(
      "/api/v1/namespaces/$NAMESPACE/pods/$PODNAME/log?timeoutSeconds=2000&watch=false",
      namespace, podName)
   req, err := http.NewRequest("GET", myurl.String(), nil)
   addAuthorization(req)

   resp, err := client.Do(req)
   if err != nil {
      panic(err.Error())
   }

   var logString string
   if resp != nil {
      data, err := ioutil.ReadAll(resp.Body)
      if err != nil {
         panic(err.Error())
      }
      logString = string(data)
      //fmt.Println(logString)

   }
   return logString

}
func getPods() []string {

   client := getClient()

   myurl := getUrl(
      "/api/v1/namespaces/$NAMESPACE/pods?timeoutSeconds=2000&watch=false",
      namespace, "")

   req, err := http.NewRequest("GET", myurl.String(), nil)
   addAuthorization(req)
   resp, err := client.Do(req)
   if err != nil {
      panic(err.Error())
   }
   if resp.StatusCode != 200 {
      err := errors.New(fmt.Sprintf("at %s, %d", resp.Status, resp.StatusCode))
      panic(err.Error())
   }

   var pods []string
   if resp != nil {
      data, err := ioutil.ReadAll(resp.Body)
      if err != nil {
         panic(err.Error())
      }
      jsonParsed, _ := gabs.ParseJSON(data)
      podItems, _ := jsonParsed.Path("items").Children()

      //      fmt.Printf("%+q\n", podItems)
      for _, pod := range podItems {
         podName := pod.Path("metadata.name").Data().(string)
         pods = append(pods, podName)
      }
   }
   return pods
}
func addAuthorization(req *http.Request) {
   req.Header.Add("Authorization", "Bearer 425BGrYBfjR8Bud7gpfiIPO9YPoOUoSDaMVcWmylyyk")
}
func getUrl(query string, namespace string, podname string) *url.URL {
   apiUri := "https://kubernetes.default.svc"
   uriString := apiUri + query
   uriString = strings.Replace(uriString,"$NAMESPACE", namespace, -1)
   uriString = strings.Replace(uriString,"$PODNAME", podname, -1)
   uri, _ := url.Parse(uriString)
   return uri
}
func getClient() (*http.Client) {
   client := &http.Client{
      Transport: &http.Transport{
         TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
         },
      },
   }
   return client
}