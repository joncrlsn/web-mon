# web-mon
A small web application monitor written in GoLang (AKA Go).  This is still in progress.  This was created for monitoring the HTTP response times on a Java application server and take thread dumps when the response times are too slow, but feel free to fork it to work for your project, or submit a pull request with generally useful changes.  

### features
* configure runtime via an external config file
* monitor as many URLs as you wish
* send an email to one or more email addresses when response time is slow or no response

### get started
* run "web-mon -g > my.config" to create an example configuration file
* update my.config with URLs you wish to monitor
* configure the email information
* test the email configuration with "web-mon --test-mail --config=my.config"
* use the default configuration settings
* run "web-mon --config=my.config" 

### todo
* Externalize the hard-coded variables in main.go to a config file
* Add the normal flags (--help, --version, --verbose, etc)
