//
// Copyright (c) 2015 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var propertySplittingRegex = regexp.MustCompile(`\s*=\s*`)
var commaSplittingRegex = regexp.MustCompile(`\s*,\s*`)

func intValue(props map[string]string, name string) (int, bool) {
	if value, ok := props[name]; ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid integer config value for %s: %s \n", name, value)
			return 0, false
		}
		return intValue, true
	}

	return 0, false
}

func boolValue(props map[string]string, name string) (bool, bool) {
	if value, ok := props[name]; ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid bool config value for %s: %s \n", name, value)
			return false, false
		}
		return boolValue, true
	}
	return false, false
}

// processConfigFile reads the properties in the given file and assigns them to global variables
func processConfigFile(fileName string) {
	if verbose {
		fmt.Println("Processing config file:", fileName)
	}
	props, err := _readPropertiesFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file "+fileName+":", err)
		os.Exit(1)
	}
	_processConfig(props)
}

func _processConfig(props map[string]string) {

	var intVal int
	var strVal string
	var boolVal bool
	var ok bool

	if boolVal, ok = boolValue(props, "verbose"); ok {
		verbose = boolVal
		if verbose {
			fmt.Println("verbose:", verbose)
		}
	}
	if intVal, ok = intValue(props, "maxResponseTimeInSeconds"); ok {
		maxResponseTime = time.Duration(intVal) * time.Second
		fmt.Println("maxResponseTime:", maxResponseTime)
	}
	if intVal, ok = intValue(props, "monitorIntervalInMinutes"); ok {
		monitorInterval = time.Duration(intVal) * time.Minute
		fmt.Println("monitorInterval:", monitorInterval)
	}
	if intVal, ok = intValue(props, "disableIntervalInMinutes"); ok {
		disableInterval = time.Duration(intVal) * time.Minute
		fmt.Println("disableInterval:", disableInterval)
	}
	if intVal, ok = intValue(props, "logIntervalInMinutes"); ok {
		logInterval = time.Duration(intVal) * time.Minute
		fmt.Println("logInterval:", logInterval)
	}
	if strVal, ok = props["shellCommand"]; ok {
		shellCommand = strVal
		fmt.Println("shellCommand:", shellCommand)
	}
	if strVal, ok = props["mailHost"]; ok {
		mailHost = strVal
		fmt.Println("mailHost:", mailHost)
	}
	if intVal, ok = intValue(props, "mailPort"); ok {
		mailPort = intVal
		fmt.Println("mailPort:", mailPort)
	}
	if strVal, ok = props["mailUsername"]; ok {
		mailUsername = strVal
		fmt.Println("mailUsername:", mailUsername)
	}
	if strVal, ok = props["mailPassword"]; ok {
		mailPassword = strVal
		fmt.Println("mailPassword: *******")
	}
	if strVal, ok = props["mailFrom"]; ok {
		mailFrom = strVal
		fmt.Println("mailFrom:", mailFrom)
	}
	if strVal, ok = props["mailTo"]; ok {
		mailTo = propertySplittingRegex.Split(strVal, -1)
		fmt.Println("mailTo:", mailTo)
	}

	//
	// Read the monitor target values.  They must be sequential like this:
	//   monitor.target1 = abc-xyz, abc-xyz.acme.com, root, joe, secret
	//   monitor.target2 = def-xyz, def-xyz.acme.com, root, joe, secret
	//   ...
	//

	targets = []Target{}
	i := 0
	for {
		i++
		if strVal, ok = props["monitor.target"+strconv.Itoa(i)]; ok {
			tgt := commaSplittingRegex.Split(strVal, 5)
			formatValid := len(tgt) > 2
			if len(tgt) > 1 {
				if strings.HasPrefix(tgt[1], "http") {
					target := Target{host: tgt[0], url: tgt[1]}
					if len(tgt) > 2 {
						target.user = tgt[2]
					}
					if len(tgt) > 3 {
						target.password = tgt[3]
					}
					targets = append(targets, target)
				} else {
					formatValid = false
				}
			}
			if !formatValid {
				fmt.Fprintln(os.Stderr, "Invalid target value:", strVal)
				fmt.Fprintln(os.Stderr, "(target value must have 2 or more comma-separated values: <host>, <url>, <httpUser>, <httpPassword>)", strVal)
			}
		} else {
			break // Assume there are no more URLs to monitor
		}

	}
}

// generateConfigurationFile prints an example configuration file to standard output
func generateConfigurationFile() {
	fmt.Println(`# web-mon configuration file.  Uncomment the values you change:
# ======================
# Monitor configuration
# ======================

# monitor.target1 = <host1>, <url1>, <httpUser>, <httpPassword>
# monitor.target2 = <host2>, <url2>, <httpUser>, <httpPassword>
# monitor.target3 = <host3>, <url3>, <httpUser>, <httpPassword>

# This is the threshold for triggering an alert.  Response times over this value create an alert
# maxResponseTimeInSeconds    = 60

# The number of minutes between monitor attempts
# monitorIntervalInMinutes    = 3

# The number of minutes monitoring will be disabled after an alert occurs
# disableIntervalInMinutes    = 60

# The number of minutes between stats logging
# logIntervalInMinutes        = 60

# A command to be executed when an alert fires
# e.g. ssh to the host and dump threads
# The hostname and process owner are passed as the arguments
# shellCommand                =

# verbose = false

# ===================
# Mail configuration
# ===================

# mailHost     = localhost
# mailPort     = 25
# mailUsername = 
# mailPassword = 

# An email address to be used as the "from" address in alert emails
# mailFrom = 

# A comma-separated list of email addresses that will receive alert emails
# mailTo = 
`)
}

// readLinesChannel reads a text file line by line into a channel.
func _readLinesChannel(filePath string) (<-chan string, error) {
	c := make(chan string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(c)
	}()
	return c, nil
}

// readPropertiesFile reads name-value pairs from a properties file
func _readPropertiesFile(fileName string) (map[string]string, error) {
	c, err := _readLinesChannel(fileName)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]string)
	for line := range c {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// Ignore this line
		} else if len(line) == 0 {
			// Ignore this line
		} else {
			parts := propertySplittingRegex.Split(line, 2)
			properties[parts[0]] = parts[1]
		}
	}

	return properties, nil
}
