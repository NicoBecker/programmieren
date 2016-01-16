package main

import (
	"log"
	"net/http"
	"os/exec"
	"io"
	"io/ioutil"
	//"bufio"
	"os"
	"fmt"
	//"time"
	"encoding/xml"
	"strconv"
	"strings"
)

var applist Applist
var running Applist
var runningHTML = make([]string, 0)
var runningProc = make([]*exec.Cmd, 0)
var path string
var appnames = make([]string,0)
var appstarts = make([]string,0)
var appstops = make([]string,0)
var outputs = make([]Output,0)
var stdinPipes = make([]io.WriteCloser, 0)
var stdOutput = make([]string, 0)

const tablestart string = "<table><tr><th>Applikation</th><th>Startbefehl</th><th>Stopbefehl</th><th>Kill</th><th>Output</th></tr>"
const tablestop string = "</table>"

type Application struct {
		App 	string 	`xml:"app,attr"`
		Start 	string	
		Stop 	string	
		Running bool
	}
	
type Applist struct {
		XMLName xml.Name `xml:"Applications"`
		Application []Application
	}
type Output struct {
		ID int
		Text string
}

func xmlRead(path string){
	file,err := os.Open(path)
	if err != nil {
		fmt.Println("Cannot read File "+path)
	}
	data,_ := ioutil.ReadAll(file)
	applist.Application = nil
	
	
	xml.Unmarshal([]byte(data), &applist)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println("App 0 Startbefehl: "+applist.Application[0].Start)
	appnames = nil
	appstarts = nil
	appstops = nil
	
	for appid := range applist.Application {
		appnames = append(appnames,applist.Application[appid].App)
		appstarts= append(appstarts,applist.Application[appid].Start)
		appstops= append(appstops,applist.Application[appid].Stop)
		applist.Application[appid].Running=false

	}
}

func printApplist() (string){
	var s string = ""
	var a = applist
	for appid := range a.Application{
		s = s+"<tr><td>"+a.Application[appid].App+"</td><td><form action='/start/?id="+strconv.Itoa(appid)+
		"' method='post'><input type='submit' value='Start'/></form></td>"
	}
	return s;
}

func generateStopButtons(){
	runningHTML = nil
	var a = running
	for appid := range a.Application{
		if running.Application[appid].Running==true{
			s:="<tr><td>"+a.Application[appid].App+"</td>"+
			"<td><form action='/stop/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Stop'/></form></td>"+
			"<td><form action='/kill/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Kill'/></form></td>"+
			"<td><form action='/output/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Output'/></form></td></tr>"
			runningHTML=append(runningHTML,s)
		} else {
			s:="<tr><td>"+a.Application[appid].App+"</td><td></td></td></td>"+
			"<td><form action='/output/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Output'/></form></td></tr>"		
			runningHTML=append(runningHTML,s)
		}
		
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("mainhandlerstart")
	path := "..\\src\\config\\config.xml"
	fmt.Println(path)
	xmlRead(path)
	//applisttoprint := printApplist()
	availabletable := tablestart+printApplist()+tablestop
	generateStopButtons()
	runningtable:= tablestart
	for id:= range runningHTML{
		runningtable = runningtable+runningHTML[id]
	}
	runningtable= runningtable+tablestop
	responseString := "<html><body><h1>Observer</h1><p><h2>Verfügbare Apps</h2>" + availabletable +"</p><p><h2>Laufende Apps</h2>"+ runningtable+"</p></body></html>"
	//fmt.Println(responseString)
	w.Write([]byte(responseString))
}

func startHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { fmt.Println(err)}
	split := strings.Split(appstarts[id]," ")
	cmd := exec.Command(split[0],split[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	} else {
		stdinPipes = append(stdinPipes, stdin)
	}

	//stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	
	err = cmd.Start()
	
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("1")
	//reader := bufio.NewReader(stdout)
	fmt.Println("2")
	//line, _, _ := reader.ReadLine()
	fmt.Println("3")
	//stdOutput = append(stdOutput, string(line))
	fmt.Println("4")
	//responseString := "<html><header><h1>startHandler</h1></header><body>" + outputs[0].Text + "</body></html>"
	running.Application = append(running.Application,applist.Application[id])
	runningProc = append(runningProc,cmd)
	running.Application[id].Running=true
	//w.Write([]byte(responseString))
	mainHandler(w, r)
}
func stopHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { fmt.Println(err)}
	split := strings.Split(running.Application[id].Stop," ")
	cmd := exec.Command(split[0],split[1:]...)
	stdout, err := cmd.Output()
	if err != nil { fmt.Println(err)}
	var output Output
	output.ID = id
	output.Text = string(stdout)
	outputs = append(outputs,output)
	running.Application[id].Running=false
	mainHandler(w, r)
}

func killHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { fmt.Println(err)}
	runningProc[id].Process.Kill()
	running.Application[id].Running=false
	mainHandler(w, r)
}

func outputHandler(w http.ResponseWriter, r *http.Request){
	
}

func main(){
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/start/", startHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/kill/", killHandler)
	http.HandleFunc("/output/", outputHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
