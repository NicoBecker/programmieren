package main

import (
	"log"
	"net/http"
	"os/exec"
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"time"
	"encoding/xml"
	"strconv"
	"strings"
)

/*

Erklärung der Variablen:

applist: Liste der Verfügbaren Anwendungen -> entspricht dem eingelesenen XML
running: Liste der Anwendungen, die gestartet wurden. Auf der Webseite in der zweiten Tabelle ersichtlich
runningHTML: Liste des HTML-Codes, mit dem die running-Tabelle erstellt wird.
runningProc: Liste der Prozesse. Hierbei handelt es sich um den realen Anwendungsprozess. Verwendet für Kill-Command
path: Pfad zur XML Datei
stdinPipes: Liste der stdin-Pipes der Anwendungen
stdoutPipes: Liste der stdout-Pipes der Anwendungen

*/

var applist Applist
var running Applist
var runningHTML = make([]string, 0)
var runningProc = make([]*exec.Cmd, 0)
var path string
var stdinPipes = make([]io.WriteCloser, 0)
var stdoutPipes = make([]io.ReadCloser, 0)

/*
	Die Folgenden Konstanten werden für die Tabellenerzeugung benutzt, siehe mainHandler
*/

const apptablestart string = "<table><tr><th>Applikation</th><th>Startbefehl</th></tr>"
const apptablestop string = "</table>"
const runtablestart string = "<table><tr><th>Applikation</th><th>Stopbefehl</th><th>Kill</th><th>Output</th><th>Autorestart</th><th>Autorestart Counter</th></tr>"
const runtablestop string = "</table>"


/*
	Die folgenden Structs sind die Grundlage für die Listen. Sie werden auf Basis des XML-Files befüllt bzw. genutzt.
*/
type Application struct {
		App 	string 	`xml:"app,attr"`
		Start 	string	
		Stop 	string	
		Running bool
		Autorestart bool
		Counter int
	}
	
type Applist struct {
		XMLName xml.Name `xml:"Applications"`
		Application []Application
	}

/*
	XMLRead iniiert die Liste "applist", indem es die Information aus dem XML-File liest. dies geschieht mit "unmarshal"
*/

func xmlRead(path string){
	file,err := os.Open(path)
	if err != nil {
		fmt.Println("Cannot read File "+path)
		fmt.Println("a")
	}
	data,_ := ioutil.ReadAll(file)
	applist.Application = nil
	xml.Unmarshal([]byte(data), &applist)
	if err != nil {
		fmt.Printf("error: %v", err)
		fmt.Println("b")
		return
	}
	
	for appid := range applist.Application {
		applist.Application[appid].Running=false
		applist.Application[appid].Autorestart=false
	}
}

/*
	printApplist bildet zunächst die Liste "applist" auf der Homepage ab, indem es eine Tabelle erzeugt
*/

func printApplist() (string){
	var s string = ""
	var a = applist
	for appid := range a.Application{
		s = s+"<tr><td>"+a.Application[appid].App+"</td><td><form action='/start/?id="+strconv.Itoa(appid)+
		"' method='post'><input type='submit' value='Start'/></form></td>"
	}
	return s;
}

/*
	generateStopButtons erzeugt die Tabelle für die laufenden Anwendungen (running).
	Dabei wird zwischen verschiedenen Flags unterschieden, die zeitgleich den Status einer Applikation widerspiegeln.
*/


func generateStopButtons(){
	runningHTML = nil
	
	var s = ""
	var a = running
	for appid := range a.Application{
		checkRunningApps(appid)
		if running.Application[appid].Running==true{
			s="<tr><td>"+a.Application[appid].App+"</td>"+
			"<td><form action='/stop/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Stop'/></form></td>"+
			"<td><form action='/kill/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Kill'/></form></td>"+
			"<td><form action='/output/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Output'/></form></td>"
			if running.Application[appid].Autorestart==true{
			s = s+ "<td><form action='/autorestart/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Autorestart ausschalten'/></form></td>"
			} else {
			s = s+ "<td><form action='/autorestart/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Autorestart einschalten'/></form></td>"
			}		
		} else {
			s="<tr><td>"+a.Application[appid].App+"</td><td></td><td></td>"+
			"<td><form action='/output/?id="+strconv.Itoa(appid)+"' method='post'><input type='submit' value='Output'/></form></td><td></td>"	
		}
		s = s+"<td> Restart Counter: "+strconv.Itoa(running.Application[appid].Counter)+"</td></tr>"
		runningHTML=append(runningHTML,s)
	}
}

/*
	checkRunningApps überprüft den Status eines Prozesses. Dadurch kann nachvollzogen werden, ob eine App von einem dritten Programm 
	(z.B. Taskmanager) geschlossen wurde. Dabei werden ebenfalls die Flags der Anwendungen beachtet und gegebenenfalls die Anwendungen neugestartet 
	bzw die Flags angepasst.
	
	Die Idee stammt von Oliver Raum, die Codeumsetzung ist an Tassilo Kloos angelehnt
*/

func checkRunningApps(appid int){
	time.Sleep(50*time.Millisecond)
	channel := make(chan error,1)
	go func(){
		channel <- runningProc[appid].Wait()
	}()
	select {
		case err := <-channel:
			if err != nil {
				err = nil
			}
		case <-time.After(100000000): 
	}
	func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(running.Application[appid].App+" mit ID "+strconv.Itoa(appid)+" exited = false")
				}				
			}()
			
	
	if runningProc[appid].ProcessState.Exited()==true && running.Application[appid].Autorestart==true{
		startsplit := strings.Split(running.Application[appid].Start, " ")
		cmd := exec.Command(startsplit[0],startsplit[1:]...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
			fmt.Println("c")
		} else {
			stdinPipes = append(stdinPipes, stdin)
		}
		if err != nil {
			log.Fatal(err)
			fmt.Println("d")
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			fmt.Println("e")
		}
		running.Application[appid].Running=true
		running.Application[appid].Counter++
		runningProc[appid]=cmd
		fmt.Println(running.Application[appid].App+" mit ID "+strconv.Itoa(appid)+" wurde neugestartet ")
	} else if runningProc[appid].ProcessState.Exited()==true && running.Application[appid].Autorestart==false{
		running.Application[appid].Running=false
		fmt.Println(running.Application[appid].App+" mit ID "+strconv.Itoa(appid)+" exited = true und autorestart false")
	} 		
	}()
}

/*
	Der mainHandler ist der Kern der Anwendung. Am ende der anderen Handler wird dieser aufgerufen, um die Webseite zu generieren.
*/

func mainHandler(w http.ResponseWriter, r *http.Request){
	//path := "..\\src\\config\\config.xml"
	path := "C:\\Users\\Nico Becker\\workspace\\programmieren\\src\\config\\config.xml"
	xmlRead(path)
	//applisttoprint := printApplist()
	availabletable := apptablestart+printApplist()+apptablestop
	generateStopButtons()
	runningtable:= runtablestart
	for id:= range runningHTML{
		runningtable = runningtable+runningHTML[id]
	}
	runningtable= runningtable+runtablestop
	responseString := "<html><body><h1>Observer</h1><p><form action='localhost:8079' method='post'><input type='submit' value='Download Config'/></form></p><p><h2>Verfügbare Apps</h2>" + availabletable +"</p><p><h2>Laufende Apps</h2>"+ runningtable+"</p></body></html>"
	//fmt.Println(responseString)
	w.Write([]byte(responseString))
}

/*
	Der startHandler startet eine Anwendung, sobald der Start Knopf betätigt wurde. Dabei wird die id über die URL übergeben und erfasst.
	Nach dem start einer Anwendung wird diese parallel in running und runningProc gelistet.
*/	

func startHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { 
		fmt.Println(err) 
		fmt.Println("f")
	}
	split := strings.Split(applist.Application[id].Start, " ")
	cmd := exec.Command(split[0],split[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		fmt.Println("g")
	} else {
		stdinPipes = append(stdinPipes, stdin)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		fmt.Println("h")
	}
	
	err = cmd.Start()
	
	if err != nil {
		fmt.Println(err)
		fmt.Println("i")
	}
	
	stdoutPipes = append(stdoutPipes, stdout)
	
	running.Application = append(running.Application,applist.Application[id])
	runningProc = append(runningProc,cmd)
	running.Application[len(running.Application)-1].Running=true
	running.Application[len(running.Application)-1].Autorestart=false
	running.Application[len(running.Application)-1].Counter=0
	mainHandler(w, r)
}

/*
	Der stopHandler führt die den im XML-File hinterlegten Stopbefehl aus. Die Anwendungen werden jedoch nicht aus den Listen running und runningProc
	gelöscht, sondern beibehalten. Dadurch kann auch nach stoppen der Anwendung noch ein Zugriff auf die jeweilige Log-Instanz gewährleistet werden.
	
	Anmerkung: der Output funktioniert nicht, Erklärung siehe Funktion "outputHandler"
*/

func stopHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { fmt.Println(err)}
	
	
	if running.Application[id].Autorestart == false {
		split := strings.Split(running.Application[id].Stop," ")
		for i:=range split{
			fmt.Println(string(split[i]))
		}
		cmd := exec.Command(split[0],split[1:]...)
		err := cmd.Start
		if err != nil { 
			fmt.Println(err)
			fmt.Println("j")
		}
		running.Application[id].Running=false
	} else {
		split := strings.Split(running.Application[id].Stop," ")
		for i:=range split{
			fmt.Println(string(split[i]))
		}
		cmd := exec.Command(split[0],split[1:]...)
		err := cmd.Start()
		if err != nil { 
			fmt.Println(err)
			fmt.Println("k")
		}
		
		running.Application[id].Running=false
		time.Sleep(50*time.Millisecond)
		
		startsplit := strings.Split(running.Application[id].Start, " ")
		cmd = exec.Command(startsplit[0],startsplit[1:]...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
			fmt.Println("l")
		} else {
			stdinPipes = append(stdinPipes, stdin)
		}
		if err != nil {
			log.Fatal(err)
			fmt.Println("m")
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			fmt.Println("n")
		}
		running.Application[id].Running=true
		running.Application[id].Counter++
		runningProc[id]=cmd
		
	}
	mainHandler(w, r)
}

/*
	Der killHandler beendet eine Anwendung über die funktion Process.kill(). Der Prozess wird über die Liste runningProc erreicht, die zugehörige
	ID wird per Query übergeben. Wie auch beim stopHandler wird dabei das jeweilige Element weder aus der running, noch aus der runningProc Liste 
	gelöscht.
*/	

func killHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { 
		fmt.Println(err)
		fmt.Println("o")
	}
	
	if running.Application[id].Autorestart==false{
		runningProc[id].Process.Kill()
		running.Application[id].Running=false
	} else {
		runningProc[id].Process.Kill()
		running.Application[id].Running=false
		time.Sleep(50*time.Millisecond)
		startsplit := strings.Split(running.Application[id].Start, " ")
		cmd := exec.Command(startsplit[0],startsplit[1:]...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
			fmt.Println("p")
		} else {
			stdinPipes = append(stdinPipes, stdin)
		}
		if err != nil {
			log.Fatal(err)
			fmt.Println("q")
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			fmt.Println("r")
		}
		running.Application[id].Running=true
		running.Application[id].Counter++
		runningProc[id]=cmd
	}
	
	mainHandler(w, r)
}

/*
	Der OutputHandler soll den Output der jeweiligen Instanz einer Anwendung anzeigen können.
*/

func outputHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { 
		fmt.Println(err)
		fmt.Println("s")
	}
	err = stdoutPipes[id].Close()
	/*
	Der Output funktioniert nicht. Wir haben mehrere Varianten getestet - bei sämtlichen Beispielanwendung gab es Fehler
	oder der Server endete in einer Endlosschleife. 
	
	Der Versuch sah folgendermaßen aus: 
	
	reader := bufio.NewReader(stdoutPipes[id])
	l,_,_ := reader.ReadLine()
	n := bytes.IndexByte(l, 0)
	s := string(l[:n])
	*/
	
	responseString := "<html><body><h1>Observer</h1><p><form><input type='button' value='Zur Startseite' onClick='history.go(-1);return true;'></form></p>"+
					"<p> Output der Anwendung "+running.Application[id].App+" mit ID "+strconv.Itoa(id)+": </p><p>Der Output funktioniert nicht. Siehe Code, Funktion : outputHandler</p></body></html>"
	w.Write([]byte(responseString))
}

/*
	Der AutorestartHandler schaltet das Autorestart Flag einer Anwendung um. (siehe Knopf auf der Webseite)
*/

func autorestartHandler(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	idq := q.Get("id")
	id, err := strconv.Atoi(idq)
	if err != nil { 
		fmt.Println(err)
		fmt.Println("s")
	}
	
	if running.Application[id].Autorestart==true{
		running.Application[id].Autorestart=false
	} else {
		running.Application[id].Autorestart=true
	}
	mainHandler(w,r)
}

/*
	In der main-Funktion werden die Handler den jeweiligen Seiten zugewiesen und der Server gestartet.
*/

func main(){
	fmt.Println("Server gestartet. Browser mit localhost:8080 aufrufen")
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/start/", startHandler)
	http.HandleFunc("/stop/", stopHandler)
	http.HandleFunc("/kill/", killHandler)
	http.HandleFunc("/output/", outputHandler)
	http.HandleFunc("/autorestart/", autorestartHandler)
	log.Fatalln(http.ListenAndServe(":8080",nil))
}
