# web-mon
A simple web site and web application monitor.  Just download the binary for your system, put it into your path, and follow the easy "getting started" instructions below.  Web-mon will monitor the HTTP response times for any URL at a configurable interval.  When the response times are too slow (or an HTTP error code is detected), an email is sent. An optional shell script can also be executed to do custom things like dump threads or capture system information.

### download
[linux64](https://github.com/joncrlsn/web-mon/raw/master/bin-linux64/web-mon "Linux 64-bit version")
[osx64](https://github.com/joncrlsn/web-mon/raw/master/bin-osx64/web-mon "OSX 64-bit version")
[win64](https://github.com/joncrlsn/web-mon/raw/master/bin-win64/web-mon.exe "Windows 64-bit version")

### features
* configure settings via an external config file
* monitor as many URLs as you wish
* supports BASIC HTTP authentication if needed (configured per URL)
* alerts via email when response time is slow, detects an error, or gets no response
* when an alert occurs, an optional external shell script can be executed.  Why?  Get thread dumps or capture system information for the host
* logs statistics since the last log message (default interval is 1 hour - configurable)

### getting started
* run "web-mon --generate-config > my.config" to create an example configuration file
* update config file with URLs you wish to monitor
* use the default monitor settings (like the interval between monitor attempts) or set your own
* configure the email settings in the config file
* test the email settings with "web-mon --test-mail --config=my.config"
* run "web-mon --config=my.config" 

### todo
* Convert flag handling to go-flags (https://github.com/jessevdk/go-flags)

