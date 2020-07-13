package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/poodlenoodle42/Film_Format_Manager/movie"
	"github.com/xfrr/goffmpeg/transcoder"
)

func getMoviesInDir(dir string, lastDir string, wg *sync.WaitGroup, movies chan<- movie.Movie) {

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
			go getMoviesInDir(dir+"/"+file.Name(), file.Name(), wg, movies)
		} else { // Found a movie
			trans := new(transcoder.Transcoder)
			err = trans.Initialize(dir+"/"+file.Name(), "")
			if err != nil {
				//fmt.Println(err)
				continue
			}
			var mov movie.Movie
			mov.Format = trans.MediaFile().Metadata().Format.FormatLongName
			mov.BitRate = trans.MediaFile().Metadata().Format.BitRate
			mov.Duration = trans.MediaFile().Metadata().Format.Duration
			mov.Name = lastDir
			mov.FileName = file.Name()
			mov.Size = int(file.Size())
			mov.Path = dir + "/" + file.Name()
			mov.NumberOfStreams = len(trans.MediaFile().Metadata().Streams)
			mov.Videostream = trans.MediaFile().Metadata().Streams[0]
			movies <- mov
		}
	}

}

func printMovStats(mov movie.Movie) {
	fmt.Println(mov.Name)
	fmt.Println("\t Filename: ", mov.FileName)
	fmt.Println("\t Format: ", mov.Format)
	fmt.Println("\t Size: ", float64(mov.Size)/1000000.0, " MB")
	fmt.Println("\t Duration: ", mov.Duration)
	fmt.Println("\t Bitrate: ", mov.BitRate)
	fmt.Println("\t Streams: ", mov.NumberOfStreams)
	fmt.Println("\t Video Stream: ")
	fmt.Println("\t\t Codec: ", mov.Videostream.CodecLongName)
	fmt.Println("\t\t Resolution: ", mov.Videostream.Width, "x", mov.Videostream.Height)
}

func printAllMoviesFullfillingReq(prequsite func(movie.Movie) bool, movies <-chan movie.Movie) {
	var moviesS []movie.Movie
	for mov := range movies {
		if prequsite(mov) {
			moviesS = append(moviesS, mov)
		}
	}
	for _, mov := range moviesS {
		printMovStats(mov)
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
	go getMoviesInDir(dir, "", &wg, movies)
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
