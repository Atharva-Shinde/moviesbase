package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/atharva-shinde/moviesbase/internal/data"
	"github.com/atharva-shinde/moviesbase/internal/validator"
)

// POST: /v1/movies
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	movie := data.Movie{}
	// importance of pointers: try this err := app.readJSON(w, r, movie)
	err := app.readJSON(w, r, &movie)
	if err != nil {
		fmt.Fprint(w, err) //personal use
		return
	}
	v := validator.New()
	v.Check(movie.Title != "", "title", "missing")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
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
