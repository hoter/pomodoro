package main

import (
	"flag"
	"fmt"
	"time"
)

type Pomodoro struct {
	work, smallRest, bigRest int
}

func main() {
	// Initialize a new Pomodoro application
	work := flag.Int("w", 25, "time for working interval (default 25 minutes)")
	smallRest := flag.Int("sr", 5, "time for small rest between working intervals (default 5 minutes)")
	bigRest := flag.Int("br", 15, "time for small rest between working intervals (default 15 minutes)")
	flag.Parse()
	pomodoro := Pomodoro{work: *work, smallRest: *smallRest, bigRest: *bigRest}

	var timer *time.Timer
	for {
		for step := 0; step < 3; step++ {
			timer = time.NewTimer(time.Minute * time.Duration(pomodoro.work))
			<-timer.C
			fmt.Println("Working interval is over. Please rest.")
			timer = time.NewTimer(time.Minute * time.Duration(pomodoro.smallRest))
			<-timer.C
		}
		timer = time.NewTimer(time.Minute * time.Duration(pomodoro.work))
		<-timer.C
		fmt.Println("You work so long. Please stand up and walk.")
		timer = time.NewTimer(time.Minute * time.Duration(pomodoro.bigRest))
		<-timer.C
	}
}
