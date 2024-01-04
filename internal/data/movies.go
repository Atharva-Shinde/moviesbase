package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/atharva-shinde/moviesbase/internal/validator"
	"github.com/lib/pq"
)

// wrapper for sql connection pool
type MovieModel struct {
	DB *sql.DB
}

func NewMovieModel(db *sql.DB) MovieModel {
	return MovieModel{
		DB: db,
	}
}

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` //this field is unexported
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	// rating    int32     // `json: "rating,omitempty,string"` // json output will be string; currently 'rating' is unexported, to change this behaviour capitalise 'r'
	Runtime int32    `json:"runtime,omitempty"`
	Genres  []string `json:"genres"`
	Version int32    `json:"version"` // starts from 1 and will increment each time movie info is updated
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

func (m MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year,runtime,genres) 
	VALUES ($1,$2,$3,$4)
	RETURNING id,created_at,version
	`

	//Scan(dest...) copies the requested row’s column values and stores them into Go values (dest...)
	//The below QueryRow returns id, created_at and version. Then the Scan() function call stores them in designated go values
	err := m.DB.QueryRow(query, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	return err
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, errors.New("no record found")
	}
	query := `
	SELECT * from movies
	WHERE id=$1
	`
	movie := Movie{}
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres), // to avoid sql: Scan error on column index 5, name "genres": unsupported Scan, storing driver.Value type []uint8 into type *[]string
		&movie.Version,
	)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
