package main

import (
	"context"
	"net/http"

	"greenlight.sparkyvxcx.co/internal/data"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Conver the string "user" to a contextKey type and assign it to the user userContextKey constant.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided User struct added
// to the context.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextGetUser() method retrieves the User struct from the request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
