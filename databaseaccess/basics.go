package databaseaccess

import (
	"database/sql"
	"fmt"

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

//DeleteTable deletes a table from a database
func DeleteTable(name string, db *sql.DB) error {
	sqlStmt := "DROP TABLE IF EXISTS \"" + name + "\";"
	_, err := db.Exec(sqlStmt)
	return err
}

//AddMovie adds a given movie to a given table in a given database
func AddMovie(mov movie.Movie, table string, db *sql.DB) error {
	sqlStmt := "INSERT INTO \"" + table + `" (Name,FileName,Path,Format,ResolutionWidth,ResolutionHeight,Codec,BitRate,Duration,Size,NumberOfStreams,Status)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?);`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		mov.Name, mov.FileName, mov.Path, mov.Format, mov.Width, mov.Height,
		mov.Codec, mov.BitRate, mov.Duration, mov.Size, mov.NumberOfStreams, mov.Status)
	if err != nil {
		return err
	}
	return nil
}

//IsMovieInDB checks if a table in a Database containes a given movie by comparing name, size and FileName
func IsMovieInDB(mov movie.Movie, table string, db *sql.DB) (bool, error) {
	sqlStmt := "SELECT EXISTS (SELECT 1 FROM \"" + table + "\" WHERE Name = ? AND Size = ? AND FileName = ?);"
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	result := stmt.QueryRow(mov.Name, mov.Size, mov.FileName)
	if err != nil {
		return false, err
	}
	var res bool
	err = result.Scan(&res)
	return res, nil
}

//OtherVersion checks if the database containes a Movie with the same name but diffrent properties
//-> returns uint
//0 "File got renamed"
//1 "Removed, diffrent file available"
//2 "Removed" / Error
func OtherVersion(mov movie.Movie, table string, db *sql.DB) (uint, error) {
	sqlStmt := "SELECT Name,Size,FileName FROM \"" + table + "\" WHERE Name = ? AND Status = 0;"
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return 2, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(mov.Name)
	if err != nil {
		return 2, err
	}
	defer rows.Close()
	for rows.Next() {
		var movR movie.Movie
		err = rows.Scan(&movR.Name, &movR.Size, &movR.FileName)
		if err != nil {
			return 2, err
		}
		if movR.Name == mov.Name && movR.Size != mov.Size {
			return 1, nil
		} else if movR.Name == mov.Name && movR.Size == mov.Size && movR.FileName != mov.FileName {
			return 0, nil
		}
	}
	return 2, nil

}

//UpdateMovieStatus Updates the status and FileName of a movie
func UpdateMovieStatus(mov movie.Movie, table string, db *sql.DB) error {
	sqlStmt := fmt.Sprintf("UPDATE \"%s\" SET Status = %d,FileName = \"%s\",Path = \"%s\" WHERE Name = \"%s\" AND Size = %d", table, mov.Status, mov.FileName, mov.Path, mov.Name, mov.Size)
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

//DeleteMovie deletes a movie based on its Name, Size and FileName
func DeleteMovie(mov movie.Movie, table string, db *sql.DB) error {
	sqlStmt := fmt.Sprintf("DELETE FROM \"%s\" WHERE Name = ? AND Size = ? AND FileName = ?", table)
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(mov.Name, mov.Size, mov.FileName)
	if err != nil {
		return err
	}
	return nil
}
