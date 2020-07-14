# Film_Format_Manager
## Whats it all about?
This little programm recursivly scans a directory containing video files and saves important meta informations. This meta information can than be searched through to find movies matching certain requirements. This can be used to gather information about the quality of your media libary. For example to keep track of which movies you already upscaled using AI.
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
The following examples are for Linux but you can easily run them on Windows by replacing `./film_format_manager_main_linux_amd64` by `film_format_manager_main_windows_amd64.exe`
### Update
This operation requires access to the folder to be scanned.
```bash
./film_format_manager_main_linux_amd64 [dir] Update
```
Updates the database with information from this directory. Preserves informations about older files by setting their status. Only scans completely new files and those changed.
### fullUpdate
This operation requires access to the folder to be scanned.
```bash
./film_format_manager_main_linux_amd64 [dir] fullUpdate
```
Updates the database with information from this directory. Deletes old information and scanns all files newly.
### list
This opertation requires no access to the folder to be searched. It only opertates on the database created using *Update* or *fullUpdate*. 
```bash
./film_format_manager_main_linux_amd64 [dir] list [criteria]
```
Prints all entrys matching a given criteria.
#### resSmallerThan / resBiggerThan
```bash
./film_format_manager_main_linux_amd64 [dir] list resSmallerThan [width] [height]
```
Prints all movies having a resolution smaller (or bigger respectively) than a threshold.
#### sizeSmallerThan / sizeBiggerThan
```bash
./film_format_manager_main_linux_amd64 [dir] list sizeSmallerThan [size in mb]
```
Prints all movies having a size smaller (or bigger respectively) than a threshold.
#### nameEq
```bash
./film_format_manager_main_linux_amd64 [dir] list nameEq [name]
```
Prints all movies having exactly this name.

#### nameCont
```bash
./film_format_manager_main_linux_amd64 [dir] list nameCont [sub_name]
```
Prints all movies which name contains the given sub_name.
#### all
```bash
./film_format_manager_main_linux_amd64 [dir] list all
```
Prints all movies.
## How to build
Given you have installed the golang packages correctly just type 
```bash
go build film_format_manager_main.go
```
in the directory you cloned this repository to.
### Windows
The sqlite driver for go uses cgo o to work. Under windows you will have to install a gcc toolchain for windows as statet here https://github.com/mattn/go-sqlite3#windows . I will try to resolve this issue.
### Why is there every exectuable twice
This programm uses concurrency to speed up its task. Your network storage might not like this.
So all executables with *_asyncpreemptoff* suffix are compiled using `GODEBUG=asyncpreemptoff=1` in front of the `go build ...` comand. The hole command to build such an executable is there for:
```bash
GODEBUG=asyncpreemptoff=1 go build film_format_manager_main.go
```
