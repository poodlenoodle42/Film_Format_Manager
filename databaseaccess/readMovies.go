package databaseaccess

import (
	"database/sql"

	"github.com/poodlenoodle42/Film_Format_Manager/movie"
)

//GetAllMovies returns a slice containing all movies from the given table
func GetAllMovies(table string, db *sql.DB) ([]movie.Movie, error) {
	sqlStmt := "SELECT * FROM \"" + table + "\";"
	return constructMoviesFromDatabase(sqlStmt, db)

}

//GetAvailableMovies returns a slice containing all movies marked as available from the given table (not used)
func GetAvailableMovies(table string, db *sql.DB) ([]movie.Movie, error) {

	sqlStmt := "SELECT * FROM \"" + table + "\" WHERE Status = 0;"
	return constructMoviesFromDatabase(sqlStmt, db)

}

//GetMoviesByName returns a slice containing all movies having a certain name (not used)
func GetMoviesByName(table string, name string, db *sql.DB) ([]movie.Movie, error) {
	sqlStmt := "SELECT * FROM \"" + table + "\" WHERE Name = \"" + name + "\";"
	return constructMoviesFromDatabase(sqlStmt, db)
}

func constructMoviesFromDatabase(sqlStmt string, db *sql.DB) ([]movie.Movie, error) {
	var movies []movie.Movie
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var mov movie.Movie
		err = rows.Scan(&id, &mov.Name, &mov.FileName, &mov.Path, &mov.Format, &mov.Width, &mov.Height,
			&mov.Codec, &mov.BitRate, &mov.Duration, &mov.Size, &mov.NumberOfStreams, &mov.Status)
		if err != nil {
			return movies, err
		}
		movies = append(movies, mov)
	}
	return movies, nil
}
