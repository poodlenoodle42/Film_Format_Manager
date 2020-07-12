package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/xfrr/goffmpeg/models"

	"github.com/xfrr/goffmpeg/transcoder"
)

type movie struct {
	name            string //Equivalent to dir containing the movie
	fileName        string //Name of the file
	path            string //path to movie
	format          string
	videostream     models.Streams
	bitRate         string
	duration        string
	size            int
	numberOfStreams int
}

func getMoviesInDir(dir string, lastDir string, wg *sync.WaitGroup, movies chan<- movie) {

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
			var mov movie
			mov.format = trans.MediaFile().Metadata().Format.FormatLongName
			mov.bitRate = trans.MediaFile().Metadata().Format.BitRate
			mov.duration = trans.MediaFile().Metadata().Format.Duration
			mov.name = lastDir
			mov.fileName = file.Name()
			mov.size = int(file.Size())
			mov.path = dir + "/" + file.Name()
			mov.numberOfStreams = len(trans.MediaFile().Metadata().Streams)
			mov.videostream = trans.MediaFile().Metadata().Streams[0]
			movies <- mov
		}
	}

}

func printMovStats(mov movie) {
	fmt.Println(mov.name)
	fmt.Println("\t Filename: ", mov.fileName)
	fmt.Println("\t Format: ", mov.format)
	fmt.Println("\t Size: ", float64(mov.size)/1000000.0, " MB")
	fmt.Println("\t Duration: ", mov.duration)
	fmt.Println("\t Bitrate: ", mov.bitRate)
	fmt.Println("\t Streams: ", mov.numberOfStreams)
	fmt.Println("\t Video Stream: ")
	fmt.Println("\t\t Codec: ", mov.videostream.CodecLongName)
	fmt.Println("\t\t Resolution: ", mov.videostream.Width, "x", mov.videostream.Height)
}

func printAllMoviesFullfillingReq(prequsite func(movie) bool, movies <-chan movie) {
	var moviesS []movie
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
	if len(os.Args) < 3 {
		panic("Not enough arguments")
	}
	dir := os.Args[1]
	mode := os.Args[2]
	if !(mode == "resSmallerThan" || mode == "resBiggerThan" || mode == "sizeSmallerThan" || mode == "sizeBiggerThan" || mode == "list") {
		panic("No known mode")
	}
	if (mode == "resSmallerThan" || mode == "resBiggerThan") && len(os.Args) != 5 {
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
	if mode == "resSmallerThan" || mode == "resBiggerThan" {
		height, err = strconv.Atoi(os.Args[4])
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	movies := make(chan movie)

	var wg sync.WaitGroup
	wg.Add(1)
	go getMoviesInDir(dir, "", &wg, movies)
	// Close the channel when all goroutines are finished
	go func() {
		wg.Wait()
		close(movies)
	}()

	if mode == "resSmallerThan" {
		printAllMoviesFullfillingReq(func(mov movie) bool {
			return mov.videostream.Height*mov.videostream.Width < width*height
		}, movies)
	} else if mode == "resBiggerThan" {
		printAllMoviesFullfillingReq(func(mov movie) bool {
			return mov.videostream.Height*mov.videostream.Width > width*height
		}, movies)
	} else if mode == "sizeSmallerThan" {
		printAllMoviesFullfillingReq(func(mov movie) bool {
			return mov.size < size
		}, movies)
	} else if mode == "sizeBiggerThan" {
		printAllMoviesFullfillingReq(func(mov movie) bool {
			return mov.size > size
		}, movies)
	} else if mode == "list" {
		printAllMoviesFullfillingReq(func(mov movie) bool {
			return true
		}, movies)
	}

}
