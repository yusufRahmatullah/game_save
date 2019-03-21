package command

import (
	"bytes"
	"testing"
)

func TestRoot(t *testing.T) {
	t.Run("show help on empty command", func(t *testing.T) {
		serv := newServiceMock()
		serv.InitGitRepo("")
		root := NewRootCommand(serv)
		buffer := bytes.Buffer{}
		helpBuffer := bytes.Buffer{}
		got := root.Parse([]string{""}, &buffer)
		want := root.Parse([]string{"help"}, &helpBuffer)
		if got != want || buffer.String() == "" || buffer.String() != helpBuffer.String() {
			t.Error("Should show help on empty command")
		}
	})
}
