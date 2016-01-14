package main

import (
	"log"
	"net/http"
	"os/exec"
	"bufio"
	"os"
	"fmt"
//	"strings"
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Text: ")
	text, _ := reader.ReadString('\n')
	txtfmt := text[0:len(text)-2]
	if txtfmt == "tasklist" {
		cmd := exec.Command(txtfmt)
		stdout, _ := cmd.Output()
		responseString := "<html><header><h1>etcHandler</h1></header><body>" + string(stdout) + "</body></html>"
		w.Write([]byte(responseString))
	} else if txtfmt == "audstart"{
		cmd := exec.Command("C:\\Program Files (x86)\\Audacity\\audacity")
		stdout, err := cmd.Output()
		fmt.Println(err)
		responseString := "<html><header><h1>etcHandler</h1></header><body>" + string(stdout) + "</body></html>"
		w.Write([]byte(responseString))
	} else if txtfmt == "audstop"{
	//	cmd:= exec.Command("tasklist")
//		stdout, _ :=cmd.Output()

		responseString := "<html><header><h1>etcHandler</h1></header><body>" + txtfmt + "</body></html>"
		w.Write([]byte(responseString))
	}

	
}

func main(){
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/etc", etcHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
