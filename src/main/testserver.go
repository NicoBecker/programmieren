package main

import (
	"log"
	"net/http"
	"os/exec"
)

func mainHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "World"
	}
	responseString := "<html><body>Hello " + name + "</body></html>"
	w.Write([]byte(responseString))
}

func dataHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "Planet"
	}
	
	cmd := exec.Command("tasklist")
	stdout, _ := cmd.Output()
	responseString := "<html><body>"+string(stdout)+"</body></html>"
	w.Write([]byte(responseString))
}

func etcHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	name := q.Get("name")
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
