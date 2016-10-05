package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Pomodoro struct {
	Port               int
	Work               Work
	SmallRest, BigRest Rest
}

func (p *Pomodoro) configFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("pomodoro.html")
		t.Execute(w, p)
	} else {
		r.ParseForm()
		p.Port, _ = strconv.Atoi(r.FormValue("Port"))
		p.Work.Argument = r.FormValue("Work.Argument")
		p.Work.Description = r.FormValue("Work.Description")
		p.Work.Value, _ = strconv.Atoi(r.FormValue("Work.Value"))
		p.SmallRest.Argument = r.FormValue("SmallRest.Argument")
		p.SmallRest.Description = r.FormValue("SmallRest.Description")
		p.SmallRest.Value, _ = strconv.Atoi(r.FormValue("SmallRest.Value"))
		p.SmallRest.Message = r.FormValue("SmallRest.Message")
		p.BigRest.Argument = r.FormValue("BigRest.Argument")
		p.BigRest.Description = r.FormValue("BigRest.Description")
		p.BigRest.Value, _ = strconv.Atoi(r.FormValue("BigRest.Value"))
		p.BigRest.Message = r.FormValue("BigRest.Message")

		configJson, _ := json.Marshal(p)
		ioutil.WriteFile("config.json", configJson, 0644)
	}
}

func (p *Pomodoro) runServer() {
	if p.Port != 0 {
		http.HandleFunc("/", p.configFormHandler)
		http.ListenAndServe(":"+strconv.Itoa(p.Port), nil)
	}
}

func (p *Pomodoro) processWork() {
	timer := time.NewTimer(time.Minute * time.Duration(p.Work.Value))
	<-timer.C
}

func (p *Pomodoro) processSmallRest() {
	fmt.Println(p.SmallRest.Message)
	timer := time.NewTimer(time.Minute * time.Duration(p.SmallRest.Value))
	<-timer.C
}

func (p *Pomodoro) processBigRest() {
	fmt.Println(p.BigRest.Message)
	timer := time.NewTimer(time.Minute * time.Duration(p.BigRest.Value))
	<-timer.C
}

type Work struct {
	Argument    string
	Value       int
	Description string
}

type Rest struct {
	Argument             string
	Value                int
	Description, Message string
}

func main() {
	// Get configuration settings from the specific file
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("File config.json was not found, please check.")
		os.Exit(1)
	}
	defer configFile.Close()

	// Initialize a new Pomodoro application
	var pomodoro Pomodoro
	if err = json.NewDecoder(configFile).Decode(&pomodoro); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	work := flag.Int(pomodoro.Work.Argument, pomodoro.Work.Value, pomodoro.Work.Description)
	smallRest := flag.Int(pomodoro.SmallRest.Argument, pomodoro.SmallRest.Value, pomodoro.SmallRest.Description)
	bigRest := flag.Int(pomodoro.BigRest.Argument, pomodoro.BigRest.Value, pomodoro.BigRest.Description)
	flag.Parse()
	pomodoro.Work.Value = *work
	pomodoro.SmallRest.Value = *smallRest
	pomodoro.BigRest.Value = *bigRest

	// Allow to users change settings
	go pomodoro.runServer()

	// Working process
	for {
		for step := 0; step < 3; step++ {
			pomodoro.processWork()
			pomodoro.processSmallRest()
		}
		pomodoro.processWork()
		pomodoro.processBigRest()
	}
}
