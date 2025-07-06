package metrics

import (
	"log"
	"sync"
	"time"
)

// Metrics provides functionality for collecting and monitoring metrics
type Metrics struct {
	requestCounts     map[string]int
	requestTimings    map[string][]time.Duration
	mu                sync.RWMutex
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		requestCounts:     make(map[string]int),
		requestTimings:    make(map[string][]time.Duration),
		mu:                sync.RWMutex{},
	}
}

// IncrementRequestCount increments the count for a specific endpoint
func (m *Metrics) IncrementRequestCount(endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestCounts[endpoint]++
	log.Printf("Request count for %s: %d", endpoint, m.requestCounts[endpoint])
}

// StartTimer returns the current time for timing a request
func (m *Metrics) StartTimer(endpoint string) time.Time {
	return time.Now()
}

// StopTimer records the time taken for a request
func (m *Metrics) StopTimer(endpoint string, start time.Time) {
	duration := time.Since(start)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestTimings[endpoint] = append(m.requestTimings[endpoint], duration)
	log.Printf("Request to %s took %v", endpoint, duration)
}

// GetRequestCount returns the number of requests for a specific endpoint
func (m *Metrics) GetRequestCount(endpoint string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.requestCounts[endpoint]
}

// GetAverageRequestTime returns the average request time for a specific endpoint
func (m *Metrics) GetAverageRequestTime(endpoint string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	timings := m.requestTimings[endpoint]
	if len(timings) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, t := range timings {
		total += t
	}
	
	return total / time.Duration(len(timings))
} 