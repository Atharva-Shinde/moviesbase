package main

import (
	"fmt"
	"net/http"
)

// GET: v1/healthcheck
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("env: %v", app.config.env) // prints to the the terminal, similar to log.Print("app.config.env")
	data := map[string]string{
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJSON(w, envelope{"data": data})
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

}
