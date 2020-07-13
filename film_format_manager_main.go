package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/poodlenoodle42/Film_Format_Manager/databaseaccess"
	"github.com/poodlenoodle42/Film_Format_Manager/movie"
	"github.com/poodlenoodle42/Film_Format_Manager/scanmethods"
)

func printAllMoviesFullfillingReq(prequsite func(movie.Movie) bool, movies <-chan movie.Movie) {
	var moviesS []movie.Movie
	db, err := databaseaccess.OpenDatabase("database.sqlite")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	err = databaseaccess.CreateNewDirectoryTable("/Videos", db)
	if err != nil {
		panic(err)
	}
	for mov := range movies {
		if prequsite(mov) {
			moviesS = append(moviesS, mov)
			b, err := databaseaccess.IsMovieInDB(mov, "/Videos", db)
			if err != nil {
				panic(err)
			}
			if !b {
				err = databaseaccess.AddMovie(mov, "/Videos", db)
			}
			if err != nil {
				panic(err)
			}
		}
	}
	for _, mov := range moviesS {
		mov.Print()
	}
}

func main() {
	modes := [...]string{"resSmallerThan", "resBiggerThan", "sizeSmallerThan", "sizeBiggerThan", "list"}
	if len(os.Args) < 3 {
		panic("Not enough arguments")
	}
	dir := os.Args[1]
	mode := os.Args[2]
	if !(mode == modes[0] || mode == modes[1] || mode == modes[2] || mode == modes[3] || mode == modes[4]) {
		panic("No known mode")
	}
	if (mode == modes[0] || mode == modes[1]) && len(os.Args) != 5 {
		panic("Not enough arguments")
	}
	size := 0
	width := 0
	var err error
	if mode != "list" {
		size, err = strconv.Atoi(os.Args[3])
		width, err = strconv.Atoi(os.Args[3])
		size *= 1000
	}
	if err != nil {
		panic(err)
	}
	height := 0
	if mode == modes[0] || mode == modes[1] {
		height, err = strconv.Atoi(os.Args[4])
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	movies := make(chan movie.Movie)
	var wg sync.WaitGroup
	wg.Add(1)
	go scanmethods.GetAllMovies(dir, "", &wg, movies)
	// Close the channel when all goroutines are finished
	go func() {
		wg.Wait()
		close(movies)
	}()

	if mode == "resSmallerThan" {
		printAllMoviesFullfillingReq(func(mov movie.Movie) bool {
			return mov.Videostream.Height*mov.Videostream.Width < width*height
		}, movies)
	} else if mode == "resBiggerThan" {
		printAllMoviesFullfillingReq(func(mov movie.Movie) bool {
			return mov.Videostream.Height*mov.Videostream.Width > width*height
		}, movies)
	} else if mode == "sizeSmallerThan" {
		printAllMoviesFullfillingReq(func(mov movie.Movie) bool {
			return mov.Size < size
		}, movies)
	} else if mode == "sizeBiggerThan" {
		printAllMoviesFullfillingReq(func(mov movie.Movie) bool {
			return mov.Size > size
		}, movies)
	} else if mode == "list" {
		printAllMoviesFullfillingReq(func(mov movie.Movie) bool {
			return true
		}, movies)
	}

}
