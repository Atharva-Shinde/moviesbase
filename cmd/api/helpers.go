package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// could be a regular function rather than a method because it doen't use any dependencies from "application"
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	strId := params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid url paramter")
	}
	return id, nil
}

// converts go objects into JSON format and writes it to http response
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// use json.NewEncoder() slightly faster than json.Marshal()
	// use json.MarshalIndent() to create a prettier terminal output of json
	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
		// w.Header().Add(key, value)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(marshalData)
	return nil
}

// converts JSON from the http request body into go object values
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, input interface{}) error {
	// set the maximum size of the request body to 1 MB, this helps in tackling Denial of Service
	// fmt.Printf("&input type %T", &input) // *interface{}
	// fmt.Printf("input type %T", input)   // *data.Movie
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// can use json.UnMarshal() instead of json.NewDecoder
	// json.NewDecoder() is more efficient than json.UnMarshal()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // disallow unknown json fields like "rating", "budget"
	// err := dec.Decode(&input) this works as well why?
	err := dec.Decode(&input) // why doesn't DisallowUnknownFields work if I don't provide pointer to movie in movies.go: app.readJSON(w, r, movie)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{}) // curl -d '{"title": "Moana"}sldklsjdl' localhost:8080/v1/movies
	if err != io.EOF {
		return errors.New("json body should contain only one string value")
	}
	// fmt.Fprintf(w, "%+v\n", input)
	return nil
}
