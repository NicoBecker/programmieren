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

func dataHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("nam")
	if name == "" {
		name = "Planet"
	}
	responseString := "<html><body>Hello " + name + "</body></html>"
	w.Write([]byte(responseString))
}

func etcHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("nam")
	if name == "" {
		name = "System"
	}
	responseString := "<html><header><h1>etcHandler</h1></header><body>Hello " + name + "</body></html>"
	w.Write([]byte(responseString))
}

func main(){
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/etc", etcHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
