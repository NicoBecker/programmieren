package main

import (
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("nam")
	if name == "" {
		name = "World"
	}
	responseString := "<html><body>Hello " + name + "</body></html>"
	w.Write([]byte(responseString))
}

func main(){
	http.HandleFunc("/", mainHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
