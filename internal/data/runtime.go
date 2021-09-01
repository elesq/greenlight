package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Abstraction which uses the same underlying type as the Movie
// struct file runtime data (int32).
type Runtime int32

// Implements a MarshalJSON to control custom return of the stringified
// representation of the our int32 value postfixed with a mins descriptor.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// Define error the override can use in event of a parsing issue
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// remove the quotes.
	unquotedJsonValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// split string to isolate the number segment
	parts := strings.Split(unquotedJsonValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil

}
