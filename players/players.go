package players

//Player interface
type Player interface {
	IsPlaying() bool
	FilePlaying() string
	PlayFile(string) error
	StopFile() error
}
