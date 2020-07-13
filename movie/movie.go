package movie

import "github.com/xfrr/goffmpeg/models"

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
