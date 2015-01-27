//
// Copyright (c) 2015 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//
package main

import (
	"fmt"
	"net/smtp"
	"time"
)

const (
	ymdhmsFormat = "2006-01-02_1504"
)

var (
	maxResponseTime    = 60 * time.Second
	threadDumpCount    = 3
	threadDumpInterval = 8 // seconds
	monitorInterval    = 3 * time.Minute
	disableInterval    = 60 * time.Minute // wait interval after a targets alert
	mailHost           = "localhost"
	mailUsername       = ""
	mailPassword       = ""
	mailFrom           = "monitor@acme.com"
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
var handleSlowResponse = func(target Target) {
	dumpJavaThreads(target.host, target.url, threadDumpCount, threadDumpInterval)
	smtp.SendMail(mailHost, nil /*no auth*/, mailFrom, mailTo, []byte("Slow response from "+target.host))
}

// main starts a go-routine for each host and url that we are monitoring
func main() {

	// alertsChan communicates errors back from the monitoring go-routines
	alertsChan := make(chan Target)

	// Start each target monitor in a go-routine
	// When a slow response or an error occurs, a monitor send an alert to the alerts channel
	for _, target := range targets {
		go monitor(target, alertsChan)
	}

	// Keep checking the alerts channel for alerts
	for {
		select {
		case tgt := <-alertsChan:
			fmt.Printf("Slow response from %s\n", tgt.host)
		default:
			fmt.Println("No message received")
			time.Sleep(2 * time.Second)
		}
	}

}

// monitor waits for a period of time then start times a request for the given URL on a regular basis.
// If the response is too slow, dump the Java threads and send an email
func monitor(target Target, alertsChan chan<- Target) {
	fmt.Printf("monitoring %s: %s\n", target.host, target.url)
	for {
		//fmt.Println("doGet(", target.url, ")")
		err := doGet(target.url)
		if err != nil {
			// Let main process know that we've found a slow system
			alertsChan <- target
			time.Sleep(disableInterval)
		} else {
			time.Sleep(monitorInterval)
		}
	}
}
