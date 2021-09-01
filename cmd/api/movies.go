package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elesq/greenlight/internal/data"
	"github.com/elesq/greenlight/internal/validator"
)

// Handler for creating a new movie record. An http.methodPost handler,
// it uses the readJSON helper function to process the request or
// triage any error. We decode to a temporary structure and then copy
// the data to a Movie struct type to avoid the ability to directly
// decode ID ot version attributes from the payload.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// Shows a movie matching the search ID. Handler uses the helper function
// readIDParam to retrieve and parse the id from the request, creates a
// movie object and calls the app.writeJSON with the movie data enveloped
// to give the response a top level entity wrapper.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
