# web-mon
A small web application monitor written in GoLang (AKA Go).  This is still in progress.  This was created for monitoring the HTTP response times on a Java application server and take thread dumps when the response times are too slow, but feel free to fork it to work for your project, or submit a pull request with generally useful changes.  

### features
* configure runtime via an external config file
* monitor as many URLs as you wish
* send an email to one or more email addresses when response time is slow or no response
* when alert occurs, execute an optional (configurable) shell script that takes the host as parameters.

### get started
* run "web-mon --generate-config > my.config" to create an example configuration file
* update config file with URLs you wish to monitor
* configure the email information in inteh config file
* test the email configuration with "web-mon --test-mail --config=my.config"
* use the default configuration settings
* run "web-mon --config=my.config" 

### todo
* Support calling HTTP resources using BASIC authentication
* Convert flag handling to go-flags (https://github.com/jessevdk/go-flags)
