package main

import (
	"net"
	"net/http"
	"time"
)

// NewTimeoutClient Copied from http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
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

// NewTimeoutClient Copied from http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}
