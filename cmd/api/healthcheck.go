package main

import (
	"net/http"
)

// GET: v1/healthcheck
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJSON(w, 201, envelope{"data": data}, nil)
	if err != nil {
		app.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

}
