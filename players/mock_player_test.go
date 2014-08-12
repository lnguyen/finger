package players

import (
	"testing"
)

func TestMockPlayerInterface(t *testing.T) {
	var _ Player = new(MockPlayer)
}
