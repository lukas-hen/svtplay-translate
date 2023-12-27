package main

import (
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	ticker   chan time.Time // Sends a "tick" (a timestamp) at a given time interval.
	requests chan string
	wg       *sync.WaitGroup
}

func NewRateLimiter(reqPerSec int) *RateLimiter {
	var wg sync.WaitGroup
	r := &RateLimiter{
		ticker:   make(chan time.Time, reqPerSec),
		requests: make(chan string, reqPerSec),
		wg:       &wg,
	}

	// Fire and forget
	// Leaking
	go func() {
		for t := range time.Tick(1000 * time.Millisecond) {
			// Ticker burst n per sec
			for i := 0; i < reqPerSec; i++ {
				r.ticker <- t
			}
		}
	}()

	return r
}

func (r *RateLimiter) Run() {
	go func() {
		for req := range r.requests {
			<-r.ticker
			log.Printf("Sent: \"%s\"\n", req)
			r.wg.Done()
		}
	}()
}

func (r *RateLimiter) Send(req string) {
	r.requests <- req
	r.wg.Add(1)
}

func (r *RateLimiter) Wait() {
	r.wg.Wait()
}
