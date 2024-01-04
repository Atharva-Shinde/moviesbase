package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atharva-shinde/moviesbase/internal/data"
)

const version = "1.0"

type config struct {
	env  string
	port int
	db   struct {
		dsn          string
		maxOpenConns int    // in-use+idle connections
		maxIdleConns int    // idle connections
		maxIdleTime  string // eg: 400ms, 3s, 13m, 1h

	}
}

// middleware of the application
type application struct {
	config config
	logger *log.Logger
	model  data.MovieModel
}

func main() {
	// declare instance of config
	cfg := config{}

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.db.maxIdleConns, "mic", 25, "maximum idle connections in Postgresql connection pool")
	flag.IntVar(&cfg.db.maxOpenConns, "moc", 25, "maximum open connections in Postgresql connection pool")
	flag.StringVar(&cfg.db.maxIdleTime, "mit", "15m", "maximum idle time a connection can exist in Postgresql connection pool")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("moviesbase_dsn"), "PostgreSQL DSN")
	// flag.StringVar(&cfg.db.dsn, "dsn", "postgres://moviesbase:password@localhost/moviesbase", "PostgreSQL DSN")
	flag.Parse()

	lg := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		lg.Fatal(err)
	}

	// close the connection pool before main function exists
	defer db.Close()

	lg.Printf("database connection pool established")

	app := application{
		config: cfg,
		logger: lg,
		model:  data.NewMovieModel(db),
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
	err = server.ListenAndServe()
	lg.Fatal(err)

}

// returns a new sql connection pool
func openDB(cfg config) (*sql.DB, error) {
	// TODO: remove the use of hardcoded values
	//creates a new connection pool for postgres
	db, err := sql.Open("postgres", "user=moviesbase password=password dbname=moviesbase sslmode=disable")
	// db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// ctx is a context with a 5 seconds timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// checks if the connection to the database is established within 5 seconds
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
