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
		// fmt.Fprint(w, err)
		return
	}
	v := validator.New()
	data.ValidateMovie(v, &movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.model.Insert(&movie)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err)
	}
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, header)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err)
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
	err = app.writeJSON(w, 201, envelope{"movie": movie}, nil)
	if err != nil {
		app.errorResponse(w, r, 500, err)
		return
	}
}
