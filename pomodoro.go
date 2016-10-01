package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Pomodoro struct {
	Work               Work
	SmallRest, BigRest Rest
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
