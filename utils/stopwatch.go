package utils

import (
	"time"
)

// StopWatch struct.
type StopWatch struct {
	start, stop time.Time
}

// Start method.
func (s *StopWatch) Start() {
	s.start = time.Now()
	s.stop = time.Time{}
}

// Stop method.
func (s *StopWatch) Stop() {
	s.stop = time.Now()
}

// Elapsed method.
func (s *StopWatch) Elapsed() time.Duration {
	if s.start.IsZero() {
		return 0
	}

	if s.stop.IsZero() {
		return time.Since(s.start)
	}

	return s.stop.Sub(s.start)
}
