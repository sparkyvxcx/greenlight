package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mux     sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while the
			// cleanup is taking place.
			mux.Lock()

			// Loop through all clients.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mux.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the client's IP address from the request.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Lock the mutex to prevent this code from being executed concurrently.
		mux.Lock()

		// Check to see if the IP address already exists in the map. If it doesn't, then
		// initialize a new rate limiter and add the IP address and limiter to the map.
		if _, found := clients[ip]; !found {
			// Initialize a new rate limiter which allows an average of 2 requests per second,
			// with maximum of 4 requests in a single 'burst'.
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		if !clients[ip].limiter.Allow() {
			mux.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// Very importantly, unlock the mutex before calling the next handler in the chain.
		// Notice that we DON't use defer to unlock the mutex, as that would mean that the
		// mutex isn't unlocked until all the handlers downstream of this middleware have
		// also returned
		mux.Unlock()

		next.ServeHTTP(w, r)
	})
}
