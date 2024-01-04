package main

import (
	"fmt"
	"net/http"

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
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	movie, err := app.model.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

// PUT: /v1/movies/:id
// updates the movie using the ID which is accessed from the http request URL and not the BODY of the request
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	movie, err := app.model.Get(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	wantMovieData := data.Movie{}
	err = app.readJSON(w, r, &wantMovieData)
	if err != nil {
		return
	}
	movie.Title = wantMovieData.Title
	movie.Genres = wantMovieData.Genres
	movie.Runtime = wantMovieData.Runtime
	movie.Year = wantMovieData.Year

	err = app.model.Update(movie)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

// DELETE /v1/movies/:id
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	err = app.model.Delete(id)
	if err != nil {
		app.errorResponse(w, r, http.StatusNotFound, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}
