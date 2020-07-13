package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/poodlenoodle42/Film_Format_Manager/scanmethods"

	"github.com/poodlenoodle42/Film_Format_Manager/databaseaccess"
	"github.com/poodlenoodle42/Film_Format_Manager/movie"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func fullUpdate(path string) {
	db, err := databaseaccess.OpenDatabase("database.sqlite")
	checkError(err)
	err = databaseaccess.DeleteTable(path, db)
	checkError(err)
	err = databaseaccess.CreateNewDirectoryTable(path, db)
	moviesChan := make(chan movie.Movie)
	var wg sync.WaitGroup
	wg.Add(1)
	scanmethods.GetAllMovies(path, "", &wg, moviesChan)
	go func() {
		wg.Wait()
		close(moviesChan)
	}()
	var moviesS []movie.Movie
	for mov := range moviesChan {
		moviesS = append(moviesS, mov)
	}
	for _, mov := range moviesS {
		mov.Print()
		err = databaseaccess.AddMovie(mov, path, db)
		checkError(err)
	}
	fmt.Println("Full Update Success")
}

func update(path string) {
	db, err := databaseaccess.OpenDatabase("database.sqlite")
	checkError(err)
	err = databaseaccess.CreateNewDirectoryTable(path, db)
	checkError(err)
	moviesChan := make(chan movie.Movie)
	var wg sync.WaitGroup
	wg.Add(1)
	go scanmethods.GetAllNewMovies(path, "", &wg, moviesChan, db, path)
	go func() {
		wg.Wait()
		close(moviesChan)
	}()
	var moviesS []movie.Movie
	for mov := range moviesChan {
		moviesS = append(moviesS, mov)
	}
	for _, mov := range moviesS {
		err = databaseaccess.AddMovie(mov, path, db)
		checkError(err)
	}
	err = scanmethods.UpdateMovieStates(path, db)
	checkError(err)
	fmt.Println("Update Success")
}

func printAllPred(pred func(movie.Movie) bool, path string) {
	db, err := databaseaccess.OpenDatabase("database.sqlite")
	checkError(err)
	movies, err := databaseaccess.GetAllMovies(path, db)
	checkError(err)
	for _, mov := range movies {
		if pred(mov) {
			mov.Print()
		}
	}
}

func main() {
	modes := [...]string{"resSmallerThan", "resBiggerThan", "sizeSmallerThan", "sizeBiggerThan", "all", "nameEq", "nameCont"}
	commands := [...]string{"fullUpdate", "Update", "list"}
	if len(os.Args) < 3 {
		panic("Not enough arguments")
	}
	path := os.Args[1]
	command := os.Args[2]
	if command == commands[0] { // fullUpdate
		fullUpdate(path)
	} else if command == commands[1] { // Update
		update(path)
	} else if command == commands[2] { // list
		mode := os.Args[3]
		if mode == modes[0] { // resSmallerThan
			width, err := strconv.Atoi(os.Args[4])
			checkError(err)
			height, err := strconv.Atoi(os.Args[5])
			checkError(err)
			printAllPred(func(mov movie.Movie) bool {
				return mov.Height*mov.Width < width*height
			}, path)
		} else if mode == modes[1] { //resBiggerThan
			width, err := strconv.Atoi(os.Args[4])
			checkError(err)
			height, err := strconv.Atoi(os.Args[5])
			checkError(err)
			printAllPred(func(mov movie.Movie) bool {
				return mov.Height*mov.Width > width*height
			}, path)
		} else if mode == modes[2] { // sizeSmallerThan
			size, err := strconv.Atoi(os.Args[4])
			checkError(err)
			printAllPred(func(mov movie.Movie) bool {
				return mov.Size < size*1000000
			}, path)
		} else if mode == modes[3] { // sizeBiggerThan
			size, err := strconv.Atoi(os.Args[4])
			checkError(err)
			printAllPred(func(mov movie.Movie) bool {
				return mov.Size > size*1000000
			}, path)
		} else if mode == modes[4] { // all
			printAllPred(func(mov movie.Movie) bool {
				return true
			}, path)
		} else if mode == modes[5] { // nameEq
			name := os.Args[4]
			printAllPred(func(mov movie.Movie) bool {
				return mov.Name == name
			}, path)
		} else if mode == modes[6] { //nameCont
			name := os.Args[4]
			printAllPred(func(mov movie.Movie) bool {
				return strings.Contains(mov.Name, name)
			}, path)
		} else {
			fmt.Println("No known command")
		}
	} else {
		fmt.Println("No known command")
	}

}
