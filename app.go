package main

import (
	"log"
	"net/http"

	"github.com/KishorPokharel/taskque/postgres"
	"github.com/julienschmidt/httprouter"
)

type application struct {
	logger  *log.Logger
	service postgres.Service
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.POST("/api/users/register", app.handleUserRegister)

	return router
}

func (app *application) run() error {
	app.logger.Println("app running")
	return http.ListenAndServe(":3000", app.routes())
}
