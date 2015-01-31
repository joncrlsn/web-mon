# web-mon
A simple web application monitor written in Go (AKA GoLang), but there should be no need to know Go.  Just download the binary for your system and put it into your path.  Web-mon can monitor the HTTP response times for any URL at a configurable interval.  When the response times are too slow, an email is sent, and an optional shell script can be executed, but feel free to fork it to work for your project, or submit a pull request with generally useful changes.  

### download
[linux64](https://github.com/joncrlsn/web-mon/raw/master/bin-linux64/web-mon "Linux 64-bit version")
[osx64](https://github.com/joncrlsn/web-mon/raw/master/bin-osx64/web-mon "OSX 64-bit version")
[win64](https://github.com/joncrlsn/web-mon/raw/master/bin-win64/web-mon.exe "Windows 64-bit version")

### features
* configure runtime via an external config file
* monitor as many URLs as you wish
* supports BASIC HTTP authentication (configured per URL)
* alerts via email when response time is slow, error, or no response
* when an alert occurs, executes an optional external shell script that takes the hostname as a parameter.  What to do with this?  Maybe you want to get thread dumps or capture system information for that host.

### getting started
* run "web-mon --generate-config > my.config" to create an example configuration file
* update config file with URLs you wish to monitor
* configure the email information in inteh config file
* test the email configuration with "web-mon --test-mail --config=my.config"
* use the default configuration settings
* run "web-mon --config=my.config" 

### todo
* Convert flag handling to go-flags (https://github.com/jessevdk/go-flags)
