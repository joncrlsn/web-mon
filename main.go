//
// Copyright (c) 2015 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//
package main

import (
	"errors"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	ymdhmsFormat = "2006-01-02_150405"
)

var (
	version         = "0.4.2"
	verbose         = false
	maxResponseTime = 60 * time.Second
	monitorInterval = 3 * time.Minute  // interval between monitoring attempts
	disableInterval = 60 * time.Minute // monitor disabled for this interval after it alerts
	mailHost        = "localhost"
	mailPort        = 25
	mailUsername    = ""
	mailPassword    = ""
	mailFrom        = ""         // an email address
	mailTo          = []string{} // a slice of email addresses
	shellCommand    = ""         // command to run when alert is triggered
)

// This is populated via the config file
var targets = []Target{}

// Target represents a hostname and a url to be monitored
type Target struct {
	host     string
	url      string
	pidOwner string
	user     string // http BASIC auth user
	password string // http BASIC auth password
	err      error
}

// doGet is overridden when testing
var doGet = func(target Target) error {

	client := NewTimeoutClient(maxResponseTime, maxResponseTime)

	req, err := http.NewRequest("GET", target.url, nil)
	if err != nil {
		log.Printf("Error creating GET request: %s: %s", target.url, err)
		return err
	}
	if len(target.user) > 0 {
		req.SetBasicAuth(target.user, target.password)
	}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting URL: %s: %s", target.url, err)
		return err
	} else {
		defer response.Body.Close()
		if response.StatusCode > 399 {
			return errors.New("Invalid response code: " + response.Status)
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response body: %s", err)
			return err
		}
		if verbose && false { // too much for verbose... should be verbose+
			log.Printf("%s\n", string(contents))
		}
	}
	if verbose {
		log.Println("response was within time limit", target.url)
	}
	return nil
}

// handleSlowResponse is overridden when testing
var handleSlowResponse = func(target *Target) {
	msg := fmt.Sprintf("Slow or error response from %s: %s: %s", target.host, target.url, target.err)
	log.Println(msg)

	var output string

	// Optionally run the shell command specified in the config file
	if len(shellCommand) > 0 {
		log.Printf("Executing shell command: %s %s %s\n", shellCommand, target.host, target.pidOwner)
		cmd := exec.Command(shellCommand, target.host, target.pidOwner)
		bytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error running shell command:", err)
		}
		output = string(bytes)
		if verbose {
			log.Printf("Output:\n %s \n", output)
		}
	}

	// Notify configured email addresses (include the output from the shell command)
	err := sendMail(msg, fmt.Sprintf("%s \n\n %s", msg, output))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error sending mail:", err)
	}
}

// processFlags returns true if processing should continue, false otherwise
func processFlags() bool {
	var configFileName string
	var versionFlag bool
	var helpFlag bool
	var generateConfig bool
	var testMail bool

	flag.StringVarP(&configFileName, "config", "c", "", "path and name of the config file")
	flag.BoolVarP(&versionFlag, "version", "V", false, "displays version information")
	flag.BoolVarP(&verbose, "verbose", "v", false, "outputs extra information")
	flag.BoolVarP(&helpFlag, "help", "?", false, "displays usage help")
	flag.BoolVarP(&generateConfig, "generate-config", "g", false, "prints a default config file to standard output")
	flag.BoolVarP(&testMail, "test-mail", "m", false, "sends a test email to the configured mail server")
	flag.Parse()

	if versionFlag {
		fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], version)
		fmt.Fprintln(os.Stderr, "\nCopyright (c) 2015 Jon Carlson.  All rights reserved.")
		fmt.Fprintln(os.Stderr, "Use of this source code is governed by the MIT license")
		fmt.Fprintln(os.Stderr, "that can be found here: http://opensource.org/licenses/MIT")
		return false
	}

	if helpFlag {
		usage()
		return false
	}

	if generateConfig {
		generateConfigurationFile()
		return false
	}

	if len(configFileName) > 0 {
		processConfigFile(configFileName)
	}

	if testMail {
		testMailConfig()
		return false
	}

	return true
}

// main starts a go-routine for each host and url that we are monitoring
func main() {

	if !processFlags() {
		// no need to proceed
		return
	}

	// alertsChan communicates errors back from the monitoring go-routines
	alertsChan := make(chan *Target)

	// Start each target monitor in a go-routine
	// When a slow response or an error occurs, a monitor send an alert to the alerts channel
	for i, target := range targets {
		if i > 0 {
			// Spread out the monitors a bit
			time.Sleep(3 * time.Second)
		}
		go monitor(target, alertsChan)
	}

	// Keep checking the alerts channel for alerts
	for {
		select {
		case tgt := <-alertsChan:
			handleSlowResponse(tgt)
		default:
			time.Sleep(5 * time.Second)
		}
	}

}

// monitor waits for a period of time then start times a request for the given URL on a regular basis.
// If the response is too slow, dump the Java threads and send an email
func monitor(target Target, alertsChan chan<- *Target) {
	log.Printf("Monitoring %s: %s\n", target.host, target.url)
	for {
		err := doGet(target) // doGet sets the err property of target
		if err != nil {
			// Let main process know that we've found a slow system
			target.err = err
			alertsChan <- &target
			time.Sleep(disableInterval)
		} else {
			target.err = nil
			time.Sleep(monitorInterval)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s --config <config-file> \n", os.Args[0])
	fmt.Fprintln(os.Stderr, `
Program flags are:
  -?, --help            : prints a summary of the commands accepted by pgrun
  -V, --version         : prints the version of pgrun being run
  -v, --verbose         : prints extra detail about what is happening
  -c, --config          : name and path of config file
  -g, --generate-config : prints an example config file to standard output
`)
}

func testMailConfig() {
	sendMail("This is a test from web-mon", "see subject")
}
