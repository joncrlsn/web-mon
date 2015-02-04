# web-mon
A simple web site and web application monitor.  Just download the binary for your system, put it into your path, and follow the easy "getting started" instructions below.  Web-mon will monitor the HTTP response times for any URL at a configurable interval.  When the response times are too slow (or an HTTP error code is detected), an email is sent. An optional shell script can also be executed to do custom things like dump threads or capture system information.

## Download
[linux64](https://github.com/joncrlsn/web-mon/raw/master/bin-linux64/web-mon "Linux 64-bit")  
[osx64](https://github.com/joncrlsn/web-mon/raw/master/bin-osx64/web-mon "OSX 64-bit")  
[win64](https://github.com/joncrlsn/web-mon/raw/master/bin-win64/web-mon.exe "Windows 64-bit")

## Features
* configure settings via an external config file
* monitor as many URLs as you wish
* supports BASIC HTTP authentication if needed (configured per URL)
* alerts via email when response time is slow, detects an error, or gets no response
* when an alert occurs, an optional external shell script can be executed.  Why?  Get thread dumps, capture system information, or whatever you want
* logs statistics since the last stats log message (default interval is 1 hour - configurable)

## Getting Started
Create an example configuration file:

      web-mon --generate-config > my.config

Update the configuration file with URLs you wish to monitor and the email settings you wish to use, etc.  Then send a test email:

      web-mon --test-mail --config=my.config

Run the monitor:

      web-mon --config=my.config

## Example config file

    # ======================
    # Monitor configuration
    # ======================

    # host, url, httpUser (optional), httpPassword (optional)
    monitor.target1 = google, http://google.com
    monitor.target2 = mywebapi, http://example.com/mywebapi, joe@example.com, super-duper-secret

    # This is the threshold for triggering an alert.  Response times over this value create an alert
    maxResponseTimeInSeconds    = 60

    # The number of minutes between monitor attempts
    monitorIntervalInMinutes    = 3

    # The number of minutes monitoring will be disabled after an alert occurs
    disableIntervalInMinutes    = 60

    # The number of minutes between each stats log message
    logIntervalInMinutes        = 60

    # A command to be executed when an alert fires
    # eg. ssh to the host and dump threads
    # The hostname is passed as an argument
    # shellCommand                =

    # verbose prints extra data to standard out
    verbose = false

    # ===================
    # Mail configuration
    # ===================

    mailHost     = smtp.example.com
    mailPort     = 25
    mailUsername = me@example.com
    mailPassword = super-secret

    # An email address to be used as the "from" address in alert emails
    mailFrom = me@example.com

    # A comma-separated list of email addresses that will receive alert emails
    mailTo = me@example.com

### ToDo
* Convert flag handling to go-flags (https://github.com/jessevdk/go-flags)
