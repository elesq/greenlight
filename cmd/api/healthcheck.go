package main

import (
	"net/http"
)

// Writes a fixed-format json response for status, operating environment and version.
// Routine uses the envelope type to wrap the data content and give it a top level
// descriptor this removing ambiguity.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	envelopedData := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelopedData, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
