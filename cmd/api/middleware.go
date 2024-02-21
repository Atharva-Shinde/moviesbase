package main

import (
	"net/http"

	"golang.org/x/time/rate"
)

// middleware to handle rate limiting
func (app *application) rateLimit(nextRequest http.Handler) http.Handler {

	// avg 2 requests per seconds and burst of 6 requests in a single burst
	limiter := rate.NewLimiter(2, 6)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tryRequest := limiter.Allow()
		if !tryRequest {
			// give 429
			app.errorResponse(w, r, 429, "rate limit exceeded")
			return
		}
		nextRequest.ServeHTTP(w, r)
	})

}
