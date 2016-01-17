package main

import (
	"net/http"
	"net/http/httptest"
	//"os/exec"
	"testing"
	"time"
)

func TestXmlRead(t *testing.T) {
	//   XML Pfad anpassen
	xmlRead("..\\config\\config.xml")
	//Test erfolgreich, falls Liste der applist.Application mit Anwendungen gefüllt wurde
	if applist.Application == nil{
		t.Error("Fehler beim Einlesen der XML-Datei")
	}
}

func TestPrintApplist(t *testing.T) {
	s := printApplist()
	if s == ""{
		t.Error("Fehler beim Erstellen des HTML-Prints für Applist")
	}
}

func TestGenerateStopButtons(t *testing.T){
	if len(running.Application) != len(runningHTML){
		t.Error("Inkonsistenz bei GenerateStopButtons: Die Listen sind nicht gleichlang")
	}
}

func TestMainHandler(t *testing.T) {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(mainHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Error zurückgegeben wird
	if err != nil {
		t.Error("Fehler beim MainHandler!")
	}
}

func TestStartHandlerLeer(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(startHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim TestStartHandlerLeer!")
	}
}

func TestStartHandler(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(startHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	url := testServer.URL+"/?id=0"
	_, err := http.Get(url)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim StartHandler")
	}
}

func TestKillHandlerLeer(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(killHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim TestKillHandlerLeer!")
	}
}

func TestKillHandler(t *testing.T){
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(killHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung dees Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?id=0"
	//Test erfolgereich, falls keine Fehler ausgegeben
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim KillHandler!")
	}
}
func TestAutorestartHandler(t *testing.T){
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(autorestartHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung dees Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?id=0"
	//Test erfolgereich, falls keine Fehler ausgegeben
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim AutorestartHandler!")
	}
}

func TestAutoRestartHandlerLeer(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(startHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim TestAutoRestartHandlerLeer!")
	}
}

func TestOutputHandler(t *testing.T){
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(outputHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung dees Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?id=0"
	//Test erfolgereich, falls keine Fehler ausgegeben
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim OutputHandler!")
	}
}

func TestOutputHandlerLeer(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(outputHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim TestOutputHandlerLeer!")
	}
}
func TestStopHandler(t *testing.T){
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(stopHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung dees Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	//Programm-ID wird übergeben
	url := testServer.URL + "/?id=0"
	//Test erfolgereich, falls keine Fehler ausgegeben
	_, err := http.Get(url)
	if err != nil {
		t.Error("Fehler beim StopHandler!")
	}
}

func TestStopHandlerLeer(t *testing.T) {
	//Programme werden aus xml-Datei ausgelesen
	xmlRead("..\\config\\config.xml")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(stopHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//Verwendung des Pakets httptest
	testServer := httptest.NewServer(s.Handler)
	_, err := http.Get(testServer.URL)
	//Test erfolgreich, falls kein Fehler ausgegeben
	//--> erwartete Ausgabe des Programms = Fehler bei strconv
	if err != nil {
		t.Error("Fehler beim TestStopHandlerLeer!")
	}
}

func TestAttributes(t *testing.T){
	if applist.Application == nil{
		t.Error("Applist leer")
	}
	if running.Application == nil{
		t.Error("Runninglist leer")
	}
	if runningHTML == nil{
		t.Error("runningHTML leer")
	}
	if runningProc == nil{
		t.Error("runningProc leer")
	}
	/*if path == ""{
		t.Error("Pfad nicht gesetzt")
	}*/
	if stdinPipes == nil{
		t.Error("stdinpipes leer")
	}
	if stdoutPipes == nil{
		t.Error("stdoutpipes leer")
	}
}


