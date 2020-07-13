package movie

import "github.com/xfrr/goffmpeg/models"

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
