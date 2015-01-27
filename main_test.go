package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Test_monitoring overrides the doGet and handleSlowResponse functions
func Test_monitoring(t *testing.T) {

	monitorInterval = 500 * time.Millisecond
	disableInterval = 3 * time.Second

	// Override the normal doGet function
	doGet = func(url string) error {
		duration := time.Duration(rand.Float64()*70) * time.Second
		fmt.Println("test: response time: ", duration)
		if duration > 60*time.Second {
			return errors.New("error")
		}
		return nil
	}

	// Override the normal handleSlowResponse function
	handleSlowResponse = func(target Target) {
		fmt.Println("test: inside handleSlowResponse function", target.url)
	}

	targets = []Target{
		Target{host: "tst-123", url: "https://tst-123/api/Ping", pidOwner: "jcarlson"},
		Target{host: "tst-abc", url: "https://tst-abc/api/Ping", pidOwner: "jcarlson"},
	}

	// alerts communicates errors back from the monitoring go-routines
	alerts := make(chan Target)

	// start each target monitor in a go-routine
	for _, target := range targets {
		fmt.Println("test: monitoring:", target)
		go monitor(target, alerts)
	}

	fmt.Println("Started monitors")

	for {
		select {
		case tgt := <-alerts:
			fmt.Printf("Slow response from %s\n", tgt.host)
		default:
			//			fmt.Println("No message received")
			time.Sleep(1 * time.Second)
		}
	}

}
