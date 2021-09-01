package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define a new type which contains a map of validation errors
type Validator struct {
	Errors map[string]string
}

// Creates a validator instance with an empty map
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// returns true if the errors map contans no entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// adds an error to the validators error map.
func (v *Validator) AddEntry(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// adds an error to the validator errors map if the check fails and is not ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddEntry(key, message)
	}
}

// returns true if a specific value is present in a list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

// returns true if a string value matches a given regex pattern
func Matches(value string, reg *regexp.Regexp) bool {
	return reg.MatchString(value)
}

// returns true if all entries in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
