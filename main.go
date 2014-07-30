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
	"github.com/longnguyen11288/arm/omxplayer"
)

var dataDir string
var channel string

type File struct {
	Filename string `json:"filename"`
}

type Status struct {
	Playing  bool   `json:"playing"`
	Filename string `json:"filename"`
}

type Files []string

func PlayFileHandler(player *omxplayer.OmxPlayer,
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

func StopFileHandler(player *omxplayer.OmxPlayer, w http.ResponseWriter) {
	err := player.StopFile()
	if err != nil {
		fmt.Fprint(w, `{ "error": "`+err.Error()+`" }`)
		return
	}
	fmt.Fprint(w, `{ "success": "true" }`)
}

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

func ChannelHandler() string {
	return fmt.Sprintf(`{ "channel": "%s" }`, channel)
}

func StatusHandler(player *omxplayer.OmxPlayer, w http.ResponseWriter) {
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
	flag.Parse()

	os.Chdir(dataDir)
	player := omxplayer.New()
	m := martini.Classic()
	m.Map(&player)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/channel", ChannelHandler)
	m.Get("/files", FilesHandler)
	m.Get("/status", StatusHandler)
	m.Post("/playfile", PlayFileHandler)
	m.Post("/stopfile", StopFileHandler)
	m.Run()
}
