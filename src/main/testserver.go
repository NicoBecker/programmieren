package main

import (
	"log"
	"net/http"
	"os/exec"
	"bufio"
	"os"
	"fmt"
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
//	reader := bufio.NewReader(os.Stdin)
//	fmt.Print("Enter Text: ")
//	text, _ := reader.ReadString('\n')
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter Text: ")
	for scanner.Scan() {
		txtfmt := scanner.Text()
	//	txtfmt := text[0:len(text)-2]
	if txtfmt == "tasklist" {
		cmd := exec.Command(txtfmt)
		stdout, _ := cmd.Output()
		responseString := "<html><header><h1>etcHandler</h1></header><body>" + string(stdout) + "</body></html>"
		w.Write([]byte(responseString))
		fmt.Print("Enter Text: ")
	} else if txtfmt == "audstart"{
		cmdToRun := "C:\\Program Files (x86)\\Audacity\\audacity.exe"
		args := []string{"arg1"}
		procAttr := new(os.ProcAttr)
		procAttr.Files = [] *os.File{os.Stdin, os.Stdout, os.Stderr}
		if process, err := os.StartProcess(cmdToRun, args, procAttr); err != nil {
			fmt.Printf("ERROR Unable to run %s: %s", cmdToRun, err.Error())
		} else {
			fmt.Printf("%s running as pid %d", cmdToRun, process.Pid)
			//process.Kill()
		}
		responseString := "<html><header><h1>etcHandler</h1></header><body>Hello</body></html>"
		w.Write([]byte(responseString))
		fmt.Print("Enter Text: ")
	} else if txtfmt == "audstop"{
		
		

		responseString := "<html><header><h1>etcHandler</h1></header><body>" + txtfmt + "</body></html>"
		w.Write([]byte(responseString))
	}
	//fmt.Print("Enter Text: ")
	}
	
}

func main(){
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/etc", etcHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
