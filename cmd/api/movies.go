package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/atharva-shinde/moviesbase/internal/data"
)

// POST: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "creating a movie")
}

// GET: /v1/movies/<id>
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.errorResponse(w, r, 500, err)
		return
	}
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Life is beautiful",
		Genres:    []string{"history", "action", "war"},
		Runtime:   122,
		Version:   1,
	}
	err = app.writeJSON(w, envelope{"movie": movie})
	if err != nil {
		app.errorResponse(w, r, 500, err)
		return
	}
}
