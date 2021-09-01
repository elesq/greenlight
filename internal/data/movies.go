package data

import (
	"time"

	"github.com/elesq/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`                // unique ID
	CreatedAt time.Time `json:"-"`                 // original timestamp movie is added to the DB, hidden.
	Title     string    `json:"title"`             // movie title
	Year      int32     `json:"year,omitempty"`    // year of release
	Runtime   Runtime   `json:"runtime,omitempty"` // runtime in minutes
	Genres    []string  `json:"genres,omitempty"`  // slice of genres
	Version   int32     `json:"version"`           // starts at 1 and is incremented anytime the data is updated
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year > 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provied")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain one or more genres")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(movie.Genres), "genres", "must not contan duplicate values")
}
