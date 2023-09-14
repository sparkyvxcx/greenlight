package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// Declare a HTTP server using the same settings as in the main() function.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a shutdownError channel to receive any errors returned by the graceful
	// Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine to listen system signals
	go func() {
		// A quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel.
		s := <-quit

		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// Create a context with a 5-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call Shutdown() on http server, passing in the context. Shutdown() will return
		// nil if the graceful shutdown was successful, or an error (which may happen becuase)
		// of a problem closing the listeners, or because the shutdown didn't complete before
		// the listeners, or because the shutdown didn't complete before the 5-second context
		// deadline is hit). Relay this return value to the shutdownError channel.
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting app server", map[string]string{"addr": srv.Addr, "env": app.config.env})

	// Calling Shutdown() on http server will cause ListenAndServe() to immediately return a
	// http.ErrServerClosed error. So if we received this error, it mean that the graceful
	// shutdown has started. So we check specifically for this, only returning the rror if it
	// is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("app server stopped", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
