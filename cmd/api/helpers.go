package main

import (
	"encoding/json"
	"errors"
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