package movie

import (
	"fmt"

	"github.com/xfrr/goffmpeg/models"
)

type Movie struct {
	Name            string //Equivalent to dir containing the movie
	FileName        string //Name of the file
	Path            string //path to movie
	Format          string
	Videostream     models.Streams
	BitRate         string
	Duration        string
	Size            int
	NumberOfStreams int
}

func (mov *Movie) Print() {
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
