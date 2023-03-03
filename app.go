package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type application struct {
	db     *sql.DB
	logger *log.Logger
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	})

	return mux
}

func (app *application) run() error {
	app.logger.Println("app running")
	return http.ListenAndServe(":3000", app.routes())
}
