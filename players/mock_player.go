package players

//MockPlayer is mock player used for testing
type MockPlayer struct {
	filename string
}

//NewMockPlayer creates new omx player
func NewMockPlayer() *MockPlayer {
	return new(MockPlayer)
}

//IsPlaying see if object is playing
func (m *MockPlayer) IsPlaying() bool {
	if m.filename != "" {
		return true
	}
	return false
}

//FilePlaying reports back the file that is being played
func (m *MockPlayer) FilePlaying() string {
	if m.filename != "" {
		return m.filename
	}
	return ""
}

//PlayFile plays a certain file on disk and returns error if it doesnt work
func (m *MockPlayer) PlayFile(filename string) error {
	m.filename = filename
	return nil
}

//StopFile stop the file that is playing
func (m *MockPlayer) StopFile() error {
	m.filename = ""
	return nil
}
