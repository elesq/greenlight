package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

// Retrieves the "id" parameter from the current request context and converts the
// value to determine the validity of an id or return an error. httprouter parses
// a request it builds a slice that contains the interpolated url parameters names
// and values.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// This function takes an http.ResponseWriter, an http.StatusCode, an enveloped
// interface{} which contains the data and a list of http headers. To create a
// response object it marshals the data passed to create a json response object.
// The marshalling process uses the json.MarshalIndent function with no prefix
// and tabs as indents. This object is appended to with newline character. The
// response is further adapted by setting a MIMEtype header and setting the header
// for the passed statusCode. Finally the Write is called to compete the handling.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	res = append(res, '\n')

	// No longer expecting any pre-write errors at this point
	// so it is now safe to write headers
	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)

	return nil
}

// Triages the error types that are likely to occur in a public-facing API
// situation. There are a mix of outcomes depending on the particular error
// these range from handled and sanitised error messages for better user
// digestion of the actual problem, there is also panic situations as well
// as catch and dispatch for errors outwith the specific cases that are
// triaged. The method works by creating a decoder with DisallowUnknownFields
// set and mitigates against dos attacks by setting a maxBytes setting to
// allow a size up to 1MB.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// Handle max sizes for mitigating DOS attacks
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Disallow unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// Triaging stage
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Syntax errors check.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d", syntaxError.Offset)

		// Secondary syntax issue possibility.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Unmarshalling type error.
		// Occurs when trying to decode into the wrong target
		// destination type.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d", unmarshalTypeError.Offset)

		// Empty body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Unknown field
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Exceeds maximum size allowed
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// Non-nil pointer error.
		// These are caught and we panic instead of returning
		// an error.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// Generic default cases for anything else.
		default:
			return err
		}
	}

	// Call Decode() again using pointer to anon empty struct
	// if the request body has a single JSON it will return an
	// io.EOF error. If the result is not the io.EOF we can
	// return a custom error.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
