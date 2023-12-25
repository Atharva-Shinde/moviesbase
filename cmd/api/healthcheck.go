package main

import (
	"fmt"
	"net/http"
)

// GET: v1/healthcheck
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("env: %v", app.config.env) // prints to the the terminal, similar to log.Print("app.config.env")
	fmt.Fprintln(w, "status available")
	fmt.Fprintf(w, "version: %v", version)

}
