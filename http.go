//
// NewTimeoutClient and TimeoutDialer copied directly from
// http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
//
package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	cookieJar = &CookieJar{jar: make(map[string][]*http.Cookie)}
)

// NewTimeoutClient returns a client that will timeout and keep in-memory cookies
func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
		Jar: cookieJar,
	}
}

// TimeoutDialer returns a connection with both a DialTimeout and a deadline for completing
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

// CookieJar holds cookies
type CookieJar struct {
	jar map[string][]*http.Cookie
}

func (p *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if verbose {
		fmt.Printf("The URL is : %s\n", u.String())
		fmt.Printf("The cookie being set is : %s\n", cookies)
	}
	p.jar[u.Host] = cookies
}

func (p *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	if verbose {
		fmt.Printf("The URL is : %s\n", u.String())
		fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	}
	return p.jar[u.Host]
}
