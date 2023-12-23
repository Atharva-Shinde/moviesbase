package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0"

type config struct {
	env  string
	port int //:8080
}

// middleware of the application
type middleware struct {
	config config
	logger *log.Logger
}

func main() {
	// declare instance of config
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	lg := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := middleware{
		config: cfg,
		logger: lg,
	}

	// multiplexer := http.NewServeMux()
	// multiplexer.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second, // server 408 error (time out)
		WriteTimeout: 10 * time.Second,
	}

	lg.Printf("starting %v server on %v", cfg.env, server.Addr)
	err := server.ListenAndServe()
	log.Fatal(err)

}
