package omxplayer

import (
	"errors"
	"os/exec"
)

//OmxPlayer struct is player that uses omxplayer
type OmxPlayer struct {
	command  *exec.Cmd
	filename string
}

//NewOmxPlayer creates new omx player
func NewOmxPlayer() *OmxPlayer {
	return new(OmxPlayer)
}

//IsPlaying see if object is playing
func (o *OmxPlayer) IsPlaying() bool {
	if o.command != nil {
		return true
	}
	return false
}

//FilePlaying reports back the file that is being played
func (o *OmxPlayer) FilePlaying() string {
	if o.filename != "" {
		return o.filename
	}
	return ""
}

//PlayFile plays a certain file on disk and returns error if it doesnt work
func (o *OmxPlayer) PlayFile(filename string) error {
	if o.IsPlaying() {
		return errors.New("Error file is playing, please stop and try again")
	}
	o.command = exec.Command("omxplayer", "-b", "-o", "local", "--loop", filename)
	o.filename = filename
	err := o.command.Start()
	if err != nil {
		return err
	}
	return nil
}

//StopFile stop the file that is playing
func (o *OmxPlayer) StopFile() error {
	o.command.Process.Kill()
	o.filename = ""
	err := exec.Command("killall", "omxplayer.bin").Run()
	return err
}
