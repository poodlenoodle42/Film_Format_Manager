package scanmethods

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/poodlenoodle42/Film_Format_Manager/databaseaccess"

	"github.com/poodlenoodle42/Film_Format_Manager/movie"
	"github.com/xfrr/goffmpeg/transcoder"
)

//GetAllMovies returns all movies in the directory and all subdirectorys through the channel
func GetAllMovies(dir string, lastDir string, wg *sync.WaitGroup, movies chan<- movie.Movie) {

	defer wg.Done()
	fmt.Println("Checking ", dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		//fmt.Println(err)
		return
	}
	for _, file := range files {
		if file.IsDir() { //Directory, search recursivly in this directory
			wg.Add(1)
			go GetAllMovies(dir+"/"+file.Name(), file.Name(), wg, movies)
		} else { // Found a movie
			trans := new(transcoder.Transcoder)
			err = trans.Initialize(dir+"/"+file.Name(), "")
			if err != nil {
				//fmt.Println(err)
				continue
			}
			var mov movie.Movie
			mov.Format = trans.MediaFile().Metadata().Format.FormatLongName
			mov.BitRate, _ = strconv.Atoi(trans.MediaFile().Metadata().Format.BitRate)
			mov.Duration, _ = strconv.Atoi(trans.MediaFile().Metadata().Format.Duration)
			mov.Name = lastDir
			mov.FileName = file.Name()
			mov.Size = int(file.Size())
			mov.Path = dir + "/" + file.Name()
			mov.NumberOfStreams = len(trans.MediaFile().Metadata().Streams)
			mov.Videostream = trans.MediaFile().Metadata().Streams[0]
			mov.Status = 0
			movies <- mov
		}
	}

}

//GetAllNewMovies only returns all movies through the channel not known already
func GetAllNewMovies(dir string, lastDir string, wg *sync.WaitGroup, movies chan<- movie.Movie, db *sql.DB, table string) {
	defer wg.Done()
	fmt.Println("Checking ", dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		//fmt.Println(err)
		return
	}
	for _, file := range files {
		if file.IsDir() { //Directory, search recursivly in this directory
			wg.Add(1)
			go GetAllMovies(dir+"/"+file.Name(), file.Name(), wg, movies)
		} else { // Found a movie

			var mov movie.Movie
			mov.Name = lastDir
			mov.Size = int(file.Size())
			b, err := databaseaccess.IsMovieInDB(mov, table, db)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			if !b { //Movie is not known
				trans := new(transcoder.Transcoder)
				err = trans.Initialize(dir+"/"+file.Name(), "")
				if err != nil {
					//fmt.Println(err)
					continue
				}
				mov.Format = trans.MediaFile().Metadata().Format.FormatLongName
				mov.BitRate, _ = strconv.Atoi(trans.MediaFile().Metadata().Format.BitRate)
				mov.Duration, _ = strconv.Atoi(trans.MediaFile().Metadata().Format.Duration)
				mov.FileName = file.Name()
				mov.Path = dir + "/" + file.Name()
				mov.NumberOfStreams = len(trans.MediaFile().Metadata().Streams)
				mov.Videostream = trans.MediaFile().Metadata().Streams[0]
				mov.Status = 0
				movies <- mov
			}
		}
	}
}

//UpdateMovieStates goes through all movies of a table, checks if they are accesible and updates their status accordingly
func UpdateMovieStates(table string, db *sql.DB) error {
	movies, err := databaseaccess.GetAllMovies(table, db)
	if err != nil {
		return err
	}
	for _, mov := range movies {
		_, err := os.Open(mov.Path)
		if err == nil { // File is accessible update status to 0
			mov.Status = 0
			err = databaseaccess.UpdateMovieStatus(mov, table, db)
			if err != nil {
				return err
			}
		} else { //File is not accessible
			b, err := databaseaccess.OtherVersion(mov, table, db)
			if err != nil {
				return err
			}
			if b { //There is an other version update status to 1
				mov.Status = 1
				err = databaseaccess.UpdateMovieStatus(mov, table, db)
				if err != nil {
					return err
				}
			} else { //There is no other version update status to 2
				mov.Status = 2
				err = databaseaccess.UpdateMovieStatus(mov, table, db)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
