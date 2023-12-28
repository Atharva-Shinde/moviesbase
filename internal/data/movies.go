package data

import (
	"time"

	"github.com/atharva-shinde/moviesbase/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` //this field is unexported
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	rating    int32     // currently 'rating' is unexported, to change this behaviour capitalise 'r'
	Runtime   int32     `json:"runtime,omitempty,string"` //json output will be a string
	Genres    []string  `json:",omitempty"`
	Version   int32     `json:"version"` // starts from 1 and will increment each time movie info is updated
}

// can have a validator receiver method instead of this function but it will affect readability: v.data.ValidateMovie(movie)
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "missing")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
