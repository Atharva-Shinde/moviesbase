package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := app.writeJSON(w, envelope{"error": message})
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}
