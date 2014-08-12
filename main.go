package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	"github.com/longnguyen11288/finger/players"
	"github.com/longnguyen11288/finger/players/omxplayer"
)

var dataDir string
var channel string
var mock bool

//File struct to output filename
type File struct {
	Filename string `json:"filename"`
}

//Status struct to output status
type Status struct {
	Playing  bool   `json:"playing"`
	Filename string `json:"filename"`
}

//Files List of files that can be played
type Files []string

//PlayFileHandler to play file
func PlayFileHandler(player players.Player,
	w http.ResponseWriter, r *http.Request) {
	var file File
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &file)
	err = player.PlayFile(file.Filename)
	if err != nil {
		fmt.Fprint(w, `{ "error": "`+err.Error()+`" }`)
		return
	}
	fmt.Fprint(w, `{ "success": "true" }`)
}

//StopFileHandler is handler to stop playing file
func StopFileHandler(player players.Player, w http.ResponseWriter) {
	err := player.StopFile()
	if err != nil {
		fmt.Fprint(w, `{ "error": "`+err.Error()+`" }`)
		return
	}
	fmt.Fprint(w, `{ "success": "true" }`)
}

//FilesHandler list the files that can be played
func FilesHandler(w http.ResponseWriter) {
	var files Files
	osFiles, _ := ioutil.ReadDir(dataDir)
	for _, f := range osFiles {
		files = append(files, f.Name())
	}
	output, err := json.Marshal(files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(output))
}

//ChannelHandler handle the channel player is listed as
func ChannelHandler() string {
	return fmt.Sprintf(`{ "channel": "%s" }`, channel)
}

//StatusHandler handles status of player
func StatusHandler(player players.Player, w http.ResponseWriter) {
	var status Status
	status.Playing = player.IsPlaying()
	status.Filename = player.FilePlaying()
	output, err := json.Marshal(status)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(output))
}

func main() {
	flag.StringVar(&dataDir, "data-dir", ".", "Data directory for videos")
	flag.StringVar(&channel, "channel", "80", "Channel used for advertisement")
	flag.BoolVar(&mock, "mock", false, "Mock server for testing")
	flag.Parse()

	os.Chdir(dataDir)
	var player players.Player
	if mock {
		player = players.NewMockPlayer()
	} else {
		player = omxplayer.NewOmxPlayer()
	}
	m := martini.Classic()
	m.Map(player)
	m.Get("/channel", ChannelHandler)
	m.Get("/files", FilesHandler)
	m.Get("/status", StatusHandler)
	m.Post("/playfile", PlayFileHandler)
	m.Post("/stopfile", StopFileHandler)
	m.Run()
}
