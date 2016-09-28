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
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("File config.json was not found, please check.")
		os.Exit(1)
	}
	defer configFile.Close()

	var pomodoro Pomodoro
	if err = json.NewDecoder(configFile).Decode(&pomodoro); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize a new Pomodoro application
	work := flag.Int(pomodoro.Work.Argument, pomodoro.Work.Value, pomodoro.Work.Description)
	smallRest := flag.Int(pomodoro.SmallRest.Argument, pomodoro.SmallRest.Value, pomodoro.SmallRest.Description)
	bigRest := flag.Int(pomodoro.BigRest.Argument, pomodoro.BigRest.Value, pomodoro.BigRest.Description)
	flag.Parse()
	pomodoro.Work.Value = *work
	pomodoro.SmallRest.Value = *smallRest
	pomodoro.BigRest.Value = *bigRest

	var timer *time.Timer
	for {
		for step := 0; step < 3; step++ {
			timer = time.NewTimer(time.Minute * time.Duration(pomodoro.Work.Value))
			<-timer.C
			fmt.Println(pomodoro.SmallRest.Message)
			timer = time.NewTimer(time.Minute * time.Duration(pomodoro.SmallRest.Value))
			<-timer.C
		}
		timer = time.NewTimer(time.Minute * time.Duration(pomodoro.Work.Value))
		<-timer.C
		fmt.Println(pomodoro.BigRest.Message)
		timer = time.NewTimer(time.Minute * time.Duration(pomodoro.BigRest.Value))
		<-timer.C
	}
}
