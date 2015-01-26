# web-mon
A small web application monitor.  Be warned, this is still in progress.  This was created for monitoring the HTTP response times on a Java application server and take thread dumps when the response times are too slow, but feel free to fork it to work for your project.  The main structure is in the main.go file and is pretty basic.

### version 0.1

### todo
* Figure out this error:  fatal error: all goroutines are asleep - deadlock!
* Externalize the hard-coded variables in main.go to a config file
* Add the normal flags (--help, --version, --verbose, etc)
