package main

import (
	"fmt"
	"time"
)

const (
	ymdhmsFormat       = "2006-01-02_1504"
	maxResponseTime    = 60 * time.Second
	threadDumpCount    = 3
	threadDumpInterval = 8 // seconds
)

// alerts communicates errors back from the monitoring go-routines
var alerts chan *Target

var monitorInterval = 3 * time.Minute
var disableInterval = 60 * time.Minute // wait interval after a targets alert

var targets = []Target{
	Target{host: "nam-msp", url: "https://web-nam-msp.crashplan.com/api/Ping", user: "central"},
	Target{host: "nbm-msp", url: "https://web-nbm-msp.crashplanpro.com/api/Ping", user: "blue"},
}

// Target represents a hostname and a url to be monitored
type Target struct {
	host string
	url  string
	user string
}

// doGet is overridden when testing
var doGet = func(url string) error {
	client := NewTimeoutClient(maxResponseTime, maxResponseTime)
	_, err := client.Get(url)
	return err
}

// handleSlowResponse is overridden when testing
var handleSlowResponse = func(target *Target) {
	dumpJavaThreads(target.host, target.url, threadDumpCount, threadDumpInterval)
	// TODO: sendEmail(target, recipientList)
}

// main starts a go-routine for each host and url that we are monitoring
func main() {
	// start each target monitor in a go-routine
	for _, target := range targets {
		go monitor(&target)
	}

	// Wait for monitors to return alerts, and restart them
	for target := range alerts {
		fmt.Printf("slow response from %s\n", target.host)
		handleSlowResponse(target)
	}
}

// monitor waits for a period of time then start times a request for the given URL on a regular basis.
// If the response is too slow, dump the Java threads and send an email
func monitor(target *Target) {
	fmt.Printf("monitoring %s: %s\n", target.host, target.url)
	for {
		err := doGet(target.url)
		if err != nil {
			// Let main process know that we've found a slow system
			alerts <- target
			time.Sleep(disableInterval)
		} else {
			time.Sleep(monitorInterval)
		}
	}
}
