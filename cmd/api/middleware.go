package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// middleware to handle rate limiting
func (app *application) rateLimit(nextRequest http.Handler) http.Handler {
	// executed only once
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	// called everytime the middleware receives a request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enable {
			// this ensures only one request from (one or many clients) is allowed to process at a time
			// as we handle the modification of the clients map and don't want any race conditions to occur
			mu.Lock()
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.errorResponse(w, r, 500, "internal server error")
			}
			if _, clientLimiter := clients[ip]; !clientLimiter {
				clients[ip] = rate.NewLimiter(rate.Limit(app.config.limiter.rate), app.config.limiter.burst)
			}
			if !clients[ip].Allow() {
				msg := fmt.Sprintf("rate limit exceeded [%s]", ip)
				app.errorResponse(w, r, 429, msg)
				// after unsuccessful operation free the lock to allow other requests to take place
				mu.Unlock()
				return
			}

			// after successful operations free the lock to allow other requests to take place
			mu.Unlock()
		}
		nextRequest.ServeHTTP(w, r)
	})
}
