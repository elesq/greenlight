package main

import (
	"fmt"
	"net/http"
)

// Generic helper for logging an error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// Generic helper for sending JSON formatted error messages to
// the client with a given statusCode. It writes its response
// using the writeJSON helper.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{
		"error": message,
	}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// server error handler 500 series
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// not found error handler 404
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// method not allowed error handler 405
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// bad request handler
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
