# Film_Format_Manager
## Whats it all about?
This little programm recursivly scans a directory containing video files and prints out important meta information when the video matches a given requirement. This can be used to gather information about the quality of your media libary and to keep track of which movies you already upscaled using AI for example.
## Preparations
### File Structure
This Programm is build for video files but it might detect other media files too. It assumes that every movie lies in a directory named like the movie in it. So assuming the following example structure: 
```bash
└── Videos
    ├── The Lord of the Rings
    │   ├──  The Lord of the Rings: The Fellowship of the Ring
    │   │   └── Some Crazy file name.mp4
    │   ├── The Lord of the Rings: The Return of the King
    │   │   └── Some Crazy file name2.mp4
    │   └──  The Lord of the Rings: The Two Towers
    │       └── Some Crazy file name3.mp4
    └── Titanic
        └── Some Crazy file name4.mp4
```
The detected name for *Some Crazy file name.mp4* would be *The Lord of the Rings: The Fellowship of the Ring*. This is usefull because you can change the name of the underlying file or the file itself and the detected name stays constant. Using recursive searching you can also have a folder containing other folders like *The Lord of the Rings* without affecting the result.
### Dependecies
You should have *ffmpeg* and *ffprobe* installed on your system and added to the path variable. You can check whether you where succesfull by typing `ffmpeg` in your terminal.  
### Get the Application
When you dont want to build it from source you can download it at the release site.
## How to use it?
The following examples are for Linux but you can easily run them on Windows by replacing `./file_format_manager_main_linux_amd64` by `file_format_manager_main_windows_amd64.exe`

### List all movies
```bash
./file_format_manager_main_linux_amd64 [dir] list 
```
### List all movies with a resolution bigger or smaller than a certain threshhold
```bash
./file_format_manager_main_linux_amd64 [dir] resSmallerThan [width] [height] 
```
```bash
./file_format_manager_main_linux_amd64 [dir] resBiggerThan [width] [height] 
```
### List all movies bigger or smaller than a certain file size

```bash
./file_format_manager_main_linux_amd64 [dir] sizeSmallerThan [size_in_mb]
```
```bash
./file_format_manager_main_linux_amd64 [dir] sizeBiggerThan [size_in_mb] 
```
## How to build
Given you have installed the golang packages correctly just type 
```bash
go build file_format_manager_main.go
```
in the directory you cloned this repository to.
### Why is there every exectuable twice
This programm uses concurrency to speed up its task. Your network storage might not like this.
So all executables with *_asyncpreemptoff* suffix are compiled using `GODEBUG=asyncpreemptoff=1` in front of the `go build ...` comand. The hole command to build such an executable is there for:
```bash
GODEBUG=asyncpreemptoff=1 go build file_format_manager_main.go
```
