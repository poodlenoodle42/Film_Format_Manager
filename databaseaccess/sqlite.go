package databaseaccess

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/poodlenoodle42/Film_Format_Manager/movie"
)

//OpenDatabase opens a database with the given name
func OpenDatabase(name string) (*sql.DB, error) {
	return sql.Open("sqlite3", name)
}

//CreateNewDirectoryTable creates a new table in a given database with a given name.
//This table can then be used by other functions to store movies in them
func CreateNewDirectoryTable(name string, db *sql.DB) error {
	sqlStmt := "CREATE TABLE IF NOT EXISTS \"" + name + `" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		Name	TEXT NOT NULL,
		FileName	TEXT NOT NULL,
		Path	TEXT NOT NULL,
		Format	TEXT NOT NULL,
		ResolutionWidth	INTEGER,
		ResolutionHeight	INTEGER,
		Codec	TEXT NOT NULL,
		BitRate	INTEGER,
		Duration	INTEGER,
		Size	INTEGER,
		NumberOfStreams	INTEGER,
		Status INTEGER
	);`
	_, err := db.Exec(sqlStmt)
	return err
}

//AddMovie adds a given movie to a given table in a given database
func AddMovie(mov movie.Movie, table string, db *sql.DB) error {
	sqlStmt := "INSERT INTO \"" + table + `" (Name,FileName,Path,Format,ResolutionWidth,ResolutionHeight,Codec,BitRate,Duration,Size,NumberOfStreams)
	VALUES (?,?,?,?,?,?,?,?,?,?,?);`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		mov.Name, mov.FileName, mov.Path, mov.Format, mov.Videostream.CodedWidth, mov.Videostream.CodedHeight,
		mov.Videostream.CodecLongName, mov.BitRate, mov.Duration, mov.Size, mov.NumberOfStreams)
	if err != nil {
		return err
	}
	return nil
}

//IsMovieInDB checks if a table in a Database containes a given movie by comparing name and size
func IsMovieInDB(mov movie.Movie, table string, db *sql.DB) (bool, error) {
	sqlStmt := "SELECT EXISTS (SELECT 1 FROM \"" + table + "\" WHERE Name = ? AND Size = ?);"
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	result := stmt.QueryRow(mov.Name, mov.Size)
	if err != nil {
		return false, err
	}
	var res bool
	err = result.Scan(&res)
	return res, nil
}

//OtherVersion checks if the database containes a Movie with the same name but diffrent properties
//-> fullfilling the requirements for status 1 "Removed, diffrent file available"
//Depending on the result the entry gets updated
func OtherVersion(mov movie.Movie, table string, db *sql.DB) (bool, error) {
	return true, nil
}
