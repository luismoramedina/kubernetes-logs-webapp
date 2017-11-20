package main

import "github.com/Jeffail/gabs"
import "crypto/tls"
import "fmt"
import "strings"
import "net/http"
import "net/url"
import "io/ioutil"

func init() {
}

func main() {
   x := getPods()
   //fmt.Println(x)
   for _, pod := range x {
     fmt.Println(pod)
   }

}

func getPods() []string {

   var myurl *url.URL
   //apiUri := "https://kubernetes.default.svc"
   apiUri := "https://api.boae.paas.gsnetcloud.corp:8443"
   namespace := "isb-npccd-dev"
   uriString := strings.Replace(
      apiUri +
      "/api/v1/namespaces/$NAMESPACE/pods?timeoutSeconds=2000&watch=false", 
      "$NAMESPACE", namespace, -1);
   myurl, err := url.Parse(uriString)
   if err != nil {
      panic(err.Error())
   }
   
   fmt.Println("before get")

   client := &http.Client{
      Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
          InsecureSkipVerify: true,
        },
      },
   }
   req, err := http.NewRequest("GET", myurl.String(), nil)
   req.Header.Add("Authorization", "Bearer 9_Zva8CgtglW3nmSyaJv-3aZejtIpn0Ymcwb_Ep13zY")
   resp, err := client.Do(req)
   if err != nil {
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