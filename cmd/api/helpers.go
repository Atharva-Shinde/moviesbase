package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// could be a regular function rather than a method because it doen't use any dependencies from "application"
func (app *application) readIDParam(w http.ResponseWriter, r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	strId := params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid url paramter")
	}
	return id, nil
}

// converts go objects into JSON format
func (app *application) writeJSON(w http.ResponseWriter, data envelope) error {
	// use json.NewEncoder() slightly faster than json.Marshal()
	// use json.MarshalIndent() to create a prettier terminal output of json
	marshalData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalData)
	return nil
}

// converts JSON in go object values
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, input interface{}) error {
	// set the maximum size of the request body to 1 MB, this helps in tackling Denial of Service
	// fmt.Printf("&input type %T", &input) // *interface{}
	// fmt.Printf("input type %T", input)   // *data.Movie
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// can use json.UnMarshal() instead of json.NewDecoder
	// json.NewDecoder() is more efficient than json.UnMarshal()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //disallow unknown json fields like "rating", "budget"
	// importance of pointers: try this err := dec.Decode(&input)
	err := dec.Decode(input)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("json body should contain only one string value")
	}
	fmt.Fprintf(w, "%+v\n", input)
	return nil
}
