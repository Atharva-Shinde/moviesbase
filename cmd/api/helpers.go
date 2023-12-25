package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

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
