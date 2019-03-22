package command

import (
	"bytes"
	"testing"
)

const (
	testNoArg  = "no argument"
	testOneArg = "one argument"
	testArgs   = "more than one"
)

func TestAdd(t *testing.T) {
	t.Run("parse one argument", func(t *testing.T) {
		testCallInit(t, true, testOneArg, "add", "game_name")
	})

	t.Run("parse more than one arguments", func(t *testing.T) {
		testCallInit(t, false, testArgs, "add", "game", "name")
	})

	t.Run("show help on parse empty arguments", func(t *testing.T) {
		testCallInit(t, false, testNoArg, "add")
	})

	t.Run("show error if not call init", func(t *testing.T) {
		testNotCallInit(t, false, "add", "game_name")
	})
}

func TestInit(t *testing.T) {
	t.Run("parse one argument", func(t *testing.T) {
		testNotCallInit(t, true, "init", "http://test.com/test.git")
	})

	t.Run("parse more than one argument", func(t *testing.T) {
		testNotCallInit(t, false, "init", "http://test.com/test.git", "another args")
	})

	t.Run("show help on empty argument", func(t *testing.T) {
		testNotCallInit(t, false, "init")
	})
}

func TestLoad(t *testing.T) {
	t.Run("parse no argument", func(t *testing.T) {
		testCallPrepared(t, true, true, testNoArg, "load")
	})

	t.Run("parse arguments", func(t *testing.T) {
		testCallPrepared(t, false, true, testOneArg, "load", "arg1")
		testCallPrepared(t, false, true, testArgs, "load", "arg1", "arg2")
	})

	t.Run("show error if not call init and set-path", func(t *testing.T) {
		testNotCallInit(t, false, "load")
		testCallInit(t, false, "not call set-path", "load")
	})
}

func TestSave(t *testing.T) {
	t.Run("parse no argument", func(t *testing.T) {
		testCallPrepared(t, true, true, testNoArg, "save")
	})

	t.Run("parse arguments", func(t *testing.T) {
		testCallPrepared(t, false, true, testOneArg, "save", "arg1")
		testCallPrepared(t, false, true, testArgs, "save", "arg1", "arg2")
	})

	t.Run("show error if not call init and set-path", func(t *testing.T) {
		testNotCallInit(t, false, "save")
		testCallInit(t, false, "not call set-path", "save")
	})
}

func TestSetPath(t *testing.T) {
	t.Run("parse one argument", func(t *testing.T) {
		testCallPrepared(t, true, false, testOneArg, "set-path", "./game/save/path")
	})

	t.Run("parse more than one arguments", func(t *testing.T) {
		testCallPrepared(t, false, false, testArgs, "set-path", "./game/save/path", "another args")
	})

	t.Run("show help on parse empty arguments", func(t *testing.T) {
		testCallPrepared(t, false, false, testNoArg, "set-path")
	})

	t.Run("show error if not call init", func(t *testing.T) {
		testNotCallInit(t, false, "set-path")
	})
}

func TestVersion(t *testing.T) {
	t.Run("return valid version", func(t *testing.T) {
		t.Helper()
		serv := newServiceMock()
		root := NewRootCommand(serv)
		mockVersion := "0.1.0"
		root.SetVersion(mockVersion)
		var buffer bytes.Buffer
		err := root.Parse([]string{"version"}, &buffer)

		if got := buffer.String(); err != nil || got != mockVersion {
			t.Errorf("Shouldn't show error. Error: %v, got: '%s', want: '%s'", err, got, mockVersion)
		}
	})
}

func testCallInit(t *testing.T, shouldPass bool, testType string, commandAndArgs ...string) {
	t.Helper()
	serv := newServiceMock()
	serv.InitGitRepo("")
	root := NewRootCommand(serv)
	testRoot(t, root, shouldPass, testType, commandAndArgs...)
}

func testCallPrepared(t *testing.T, shouldPass bool, withSavePath bool, testType string, commandAndArgs ...string) {
	t.Helper()
	serv := newServiceMock()
	serv.InitGitRepo("")
	serv.AddConfig("game_name", "game1")
	serv.PrepareGame()
	if withSavePath {
		serv.AddConfig("save_path", "./dummy/path")
	}
	root := NewRootCommand(serv)
	testRoot(t, root, shouldPass, testType, commandAndArgs...)
}

func testNotCallInit(t *testing.T, shouldPass bool, commandAndArgs ...string) {
	t.Helper()
	serv := newServiceMock()
	root := NewRootCommand(serv)
	testRoot(t, root, shouldPass, "uninitialized", commandAndArgs...)
}

func testRoot(t *testing.T, root *RootCommand, shouldPass bool, testType string, commandAndArgs ...string) {
	t.Helper()
	var buffer bytes.Buffer
	err := root.Parse(commandAndArgs, &buffer)
	if !shouldPass && err == nil {
		t.Errorf("Should show error on %s", testType)
	} else if shouldPass && err != nil {
		t.Errorf("Shouldn't show error on %s, error: %v", testType, err)
	}
}
