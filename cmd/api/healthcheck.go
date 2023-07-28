package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// js := `{"status": "available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request at the moment", http.StatusInternalServerError)
	}

	// Append a newline to the JSON data to make it easier to view in terminal.
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
