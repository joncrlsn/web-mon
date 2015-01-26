package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
}

// Test_monitoring overrides the doGet and handleSlowResponse functions
func Test_monitoring(t *testing.T) {

	monitorInterval = 2 * time.Nanosecond
	disableInterval = 3 * time.Second

	// Override the normal doGet function
	doGet = func(url string) error {
		duration := time.Duration(rand.Float64()*300) * time.Second
		fmt.Println("test: response time: ", duration)
		if duration > 60*time.Second {
			return errors.New("error")
		}
		return nil
	}

	// Override the normal handleSlowResponse function
	handleSlowResponse = func(target *Target) {
		fmt.Println("test: inside handleSlowResponse function", target.url)
	}

	targets = []Target{
		Target{host: "tst-msp", url: "https://web-tst-msp/api/Ping", user: "central"},
		Target{host: "tst-sea", url: "https://web-tst-sea/api/Ping", user: "blue"},
	}

	// start each target monitor in a go-routine
	for _, target := range targets {
		fmt.Println("test: monitoring:", target)
		go monitor(&target)
	}

	fmt.Println("Started monitors")

	// Wait for monitors to return alerts and do something about it.
	for target := range alerts {
		fmt.Printf("test: slow response from %s\n", target.host)
		handleSlowResponse(target)
	}

}
