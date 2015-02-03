//
// Copyright (c) 2015 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//
package main

import (
	"fmt"
	"time"
)

// Stats instances hold a few statistics about HTTP requests
type Stats struct {
	StartTime         time.Time
	SampleCount       int
	TotalResponseTime time.Duration
	MaxResponseTime   time.Duration
	MinResponseTime   time.Duration
}

// Add an HTTP timing to the stats
func (s *Stats) Add(d time.Duration) {
	s.SampleCount++
	s.TotalResponseTime += d
	if d > s.MaxResponseTime {
		s.MaxResponseTime = d
	}
	if s.MinResponseTime == 0 || d < s.MinResponseTime {
		s.MinResponseTime = d
	}
}

// AvgResponseTime returns the average response time since the last call to Clear()
func (s *Stats) AvgResponseTime() time.Duration {
	return time.Duration(int64(s.TotalResponseTime) / int64(s.SampleCount))
}

// Clear the stats to start over again
func (s *Stats) Clear() {
	s.StartTime = time.Now()
	s.SampleCount = 0
	s.TotalResponseTime = time.Duration(0)
	s.MaxResponseTime = time.Duration(0)
	s.MinResponseTime = time.Duration(0)
}

// String returns a string representation of the stats
func (s *Stats) String() string {
	return fmt.Sprintf("Stats: count:%d, avgResponse:%v, maxResponse:%v, minResponse:%v", s.SampleCount, s.AvgResponseTime(), s.MaxResponseTime, s.MinResponseTime)
}
