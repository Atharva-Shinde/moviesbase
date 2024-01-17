package main

import (
	"database/sql"
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

	// !check the comment on line 86
	//wantMovieData:= data.Movie{}
	wantMovieData := struct {
		Title   *string  `json:"title"`
		Runtime *int32   `json:"runtime"`
		Year    *int32   `json:"year"`
		Genres  []string `json:"genres"` // no need to introduce a pointer, as slices can are nil by default
	}{}
	err = app.readJSON(w, r, &wantMovieData)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "bad json response")
		return
	}
	// we want to let users update values in the movie w/o the need to provide an entire json containing all the fields and their values
	// to achieve this we need to tell our code to check only for the provided fields
	// but this is what happens if we try comparing the fields(title, year etc.) with the actual data; eg: wantMovieData.Year != 0{....}
	// we are restricted to compare Year with zero! And this not a behaviour we desire
	// therefore rather than checking if the year field is zero or not we should rather validate if the year field is nil
	// to achive this, we need to convert the wantMovieData fields to pointers so that all fields can be compared to nil
	// if wantMovieData.Title != "" {
	// 	movie.Title = wantMovieData.Title
	// }
	// if wantMovieData.Genres != nil {
	// 	movie.Genres = wantMovieData.Genres
	// }
	// if wantMovieData.Runtime != 0 {
	// 	movie.Runtime = wantMovieData.Runtime
	// }
	// if wantMovieData.Year != 0 {
	// 	movie.Year = wantMovieData.Year
	// }

	// update every fields that the users provides into the movie object
	if wantMovieData.Title != nil {
		movie.Title = *wantMovieData.Title // dereference the title
	}
	if wantMovieData.Runtime != nil {
		movie.Runtime = *wantMovieData.Runtime
	}
	if wantMovieData.Year != nil {
		movie.Year = *wantMovieData.Year
	}
	if wantMovieData.Genres != nil {
		movie.Genres = wantMovieData.Genres
	}
	v := validator.New()
	data.ValidateMovie(v, movie)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// TODO: check if the new data is same as the existing data is affirmative, gracefully log it
	err = app.model.Update(movie)
	if err != nil {
		app.errorResponse(w, r, http.StatusConflict, sql.ErrNoRows)
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

type params struct {
	Title    string
	Genres   []string
	Page     int
	PageSize int
	Sort     string
	// Runtime int32
	// Year int32
}

// GET /v1/movies
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	queryableParameters := params{}
	v := validator.New()
	queryValues := r.URL.Query()
	queryableParameters.Title = app.readString(queryValues, "title")
	queryableParameters.Page = app.readInt(queryValues, "page", v)
	queryableParameters.PageSize = app.readInt(queryValues, "page_size", v)
	queryableParameters.Genres = app.readCSV(queryValues, "genres")
	queryableParameters.Sort = app.readString(queryValues, "sort")

	v.Check(queryableParameters.Page > 0, "page", "must be greater than zero")
	v.Check(queryableParameters.Page <= 10_000, "page", "must be less than ten thousand")
	v.Check(queryableParameters.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(queryableParameters.PageSize < 100, "page_size", "must be less than hundred")
	// TODO: validation check for sort

	if !v.Valid() {
		app.errorResponse(w, r, http.StatusNotAcceptable, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", queryableParameters)
}
