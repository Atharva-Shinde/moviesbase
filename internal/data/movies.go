package data

import "time"

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
