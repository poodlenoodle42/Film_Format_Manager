package movie

import (
	"fmt"
)

//Movie contains important information about a movie file
type Movie struct {
	Name            string //Equivalent to dir containing the movie
	FileName        string //Name of the file
	Path            string //path to movie
	Format          string
	Codec           string
	Width           int
	Height          int
	BitRate         int
	Duration        int
	Size            int
	NumberOfStreams int
	Status          int
}

var states = map[int]string{
	0: "Available",
	1: "Removed, diffrent file available",
	2: "Removed"}

//Print prints a Movie to the standart output
func (mov *Movie) Print() {
	fmt.Println(mov.Name)
	fmt.Println("\t Filename: ", mov.FileName)
	fmt.Println("\t State: ", states[mov.Status])
	fmt.Println("\t Format: ", mov.Format)
	fmt.Println("\t Size: ", float64(mov.Size)/1000000.0, " MB")
	fmt.Println("\t Duration: ", mov.Duration)
	fmt.Println("\t Bitrate: ", mov.BitRate)
	fmt.Println("\t Streams: ", mov.NumberOfStreams)
	fmt.Println("\t Video Stream: ")
	fmt.Println("\t\t Codec: ", mov.Codec)
	fmt.Println("\t\t Resolution: ", mov.Width, "x", mov.Height)
}
