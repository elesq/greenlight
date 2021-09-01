package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// temporary hardcoded version controller
const version = "0.1.0"

// app runtime configuration.
// settings will be read in and held.
type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {

	var cfg config

	// Read port and cli flags into the congig struct. Default port
	// number is ::6118" with a "development" env.
	flag.IntVar(&cfg.port, "port", 6118, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "environment (development|staging|production)")
	flag.Parse()

	// Define a new logger for the stdout and a date & time prefix.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// The app instance
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare the http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server and log the listening.
	logger.Printf("env %s, Listening on: %s", cfg.env, srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
