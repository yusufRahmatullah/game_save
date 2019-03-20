package command

import (
	"bytes"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("parse one argument", func(t *testing.T) {
		root := NewRootCommand(nil, nil)
		buffer := bytes.Buffer{}
		got := root.Parse([]string{"add", "game_name"}, &buffer)
		if got != nil {
			t.Error("Shouldn't show error on one argument")
		}
	})

	t.Run("parse more than one arguments", func(t *testing.T) {
		root := NewRootCommand(nil, nil)
		buffer := bytes.Buffer{}
		got := root.Parse([]string{"add", "game", "name"}, &buffer)
		if got == nil {
			t.Error("Should show error on more than one arguments")
		}
	})

	t.Run("show help on parse empty arguments", func(t *testing.T) {
		root := NewRootCommand(nil, nil)
		buffer := bytes.Buffer{}
		got := root.Parse([]string{"add"}, &buffer)
		if got == nil || buffer.String() == "" {
			t.Error("Should show error on empty arguments")
		}
	})
}
