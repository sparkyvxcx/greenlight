package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"greenlight.sparkyvxcx.co/internal/data"
	"greenlight.sparkyvxcx.co/internal/validator"
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
		if app.config.limiter.enabled {
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

			clients[ip].lastSeen = time.Now()

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
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any caches
		// that the response may vary based on the value of the Authorization header in the
		// request.
		w.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request. This will return the
		// empty string "" if there is no such header found.
		authorizationHeader := r.Header.Get("Authorization")

		// If there is no authorization header found, use the contextSetUser() helper to add the
		// AnonymousUser to the request context. Then call the next handler in the chain and return
		// without executing any of the code below.
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise expect the value of the Authorization header to be in the format "Bearer <token>".
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract authentication token from header parts
		token := headerParts[1]

		// Validate the token to ensure it is in a sensible format.
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Retrieve the user details associated with the authentication token.
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Call the contextSetUser() helper to add the user information to the request context.
		r = app.contextSetUser(r, user)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// User the contextGetUser() helper to retrieve the user information from the request context.
		user := app.contextGetUser(r)

		// If the user is anonymous, then call the authenticationRequireResponse() to inform the client
		// that they should authenticate before trying again.
		if user.IsAnonymous() {
			app.authenticateRequiredResponse(w, r)
			return
		}

		// If the user is not activated, use the inactiveAccountResponse() helper to inform them that
		// they need to activate their account.
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
