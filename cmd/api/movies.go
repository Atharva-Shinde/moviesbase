package main

import (
	"fmt"
	"net/http"
)

// POST: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "creating a movie")
}

// GET: /v1/movies/<id>
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "id of movie: %v", id)
}
