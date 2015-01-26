package main

import (
	"fmt"
	"net/smtp"
	"time"
)

const (
	ymdhmsFormat = "2006-01-02_1504"
)

// alerts communicates errors back from the monitoring go-routines
var alerts chan *Target

var (
	maxResponseTime    = 60 * time.Second
	threadDumpCount    = 3
	threadDumpInterval = 8 // seconds
	monitorInterval    = 3 * time.Minute
	disableInterval    = 60 * time.Minute // wait interval after a targets alert
	mailHost           = "localhost"
	mailUsername       = ""
	mailPassword       = ""
	mailFrom           = "joe@acme.com"
	mailTo             = []string{"joe@acme.com"}
)

var targets = []Target{
	Target{host: "xyz", url: "https://xyz.acme.com/api/Ping", pidOwner: "central"},
	Target{host: "abc", url: "https://abc.acme.com/api/Ping", pidOwner: "blue"},
}

// Target represents a hostname and a url to be monitored
type Target struct {
	host     string
	url      string
	pidOwner string
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
	smtp.SendMail(mailHost, nil /*no auth*/, mailFrom, mailTo, []byte("Slow response from "+target.host))
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
