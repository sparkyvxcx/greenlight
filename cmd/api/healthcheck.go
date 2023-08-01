package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// js := `{"status": "available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)

	env := envelop{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request at the moment", http.StatusInternalServerError)
	}
}

func healthcheckHandlerWithEncoder(w http.ResponseWriter, r *http.Request) {
	history := struct {
		v1 string
		v2 string
	}{
		v1: "2021-01-03",
		v2: "2022-06-28",
	}
	data := map[string]interface{}{
		"status":      "available",
		"environment": "dev",
		"version":     version,
		"history":     history,
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func healthcheckHandlerWithMarshal(w http.ResponseWriter, r *http.Request, beatify bool) {
	history := struct {
		v1 string
		v2 string
	}{
		v1: "2021-01-03",
		v2: "2022-06-28",
	}

	data := map[string]interface{}{
		"status":      "available",
		"environment": "dev",
		"version":     version,
		"history":     history,
	}

	w.Header().Set("Content-Type", "application/json")

	var js []byte
	var err error

	if beatify {
		js, err = json.MarshalIndent(data, "", "\t")
	} else {
		js, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
