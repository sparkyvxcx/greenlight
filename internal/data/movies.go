package data

import (
	"database/sql"
	"errors"
	"time"

	"greenlight.sparkyvxcx.co/internal/validator"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	// check `Year` field
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// check `Runtime` field
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// check `Genres` field
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MovieModel struct {
	DB *sql.DB
}

// The Insert() method accepts a pointer to a movie struct, which should contain the data for the new record
func (m MovieModel) Insert(movie *Movie) error {
	query := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	// Create a args slice containing the values for the placeholder parameters rom the movie struct. Declaring
	// this slice immediately next to our SQL query helps to make it nice and clear *what values are being uesd
	// where* in the query.
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method to execute the SQL query on our connection pool, passing in the args slice as
	// a variadic parameter and scanning the system-generated id, created_at and version values into the movie
	// struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version FROM movies WHERE id = $1`

	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5
	RETURNING version
	`

	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	// Use the QueryRow() method to execute the query, passing in the args slice as a variadic parameter and
	// and scanning the new version value into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	// Return an ErrRecordNotFound err if id is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	return nil
}
