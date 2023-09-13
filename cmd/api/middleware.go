package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A deferred function which will always be run in the event of a panic as Go
		// unwinds the stack.
		defer func() {
			// Use built-in recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				// If there was a panic, close the current connection after a response has
				// been sent.
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Initialize a new rate limiter which allows an average of 2 requests per second,
	// with maximum of 4 requests in a single 'burst'.
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
