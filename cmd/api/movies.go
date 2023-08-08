package main

import (
	"fmt"
	"net/http"

	"greenlight.sparkyvxcx.co/internal/data"
	"greenlight.sparkyvxcx.co/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Create a new movie"))

	// A anonymous struct to hold the information that exepected to be in the HTTP request body.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// fmt.Fprintf(w, "%+v\n", input)
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the client know
	// which URL they can find the newly-created resource at. We make an empty http.Header map and
	// then use the Set() method to add a new Location header, intepolating the system-generated ID
	// for the new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// movie := data.Movie{
	// 	ID:        id,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Casablanca",
	// 	Runtime:   102,
	// 	Genres:    []string{"drama", "romance", "war"},
	// 	Version:   1,
	// }

	// fmt.Fprintf(w, "Show the details of movie %d\n", id)
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
