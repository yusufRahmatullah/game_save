package repository

import (
	"path"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("copy file from a location to GameSaveRoot", func(t *testing.T) {
		rep := OSRepository{}
		srcFile := "test_copy.txt"
		createDummyFile(t, srcFile)
		ensureCloned(t, emptyRepo)
		err := rep.Copy(srcFile, GameSaveRoot)
		assertNotError(t, err)
		assertExist(t, path.Join(GameSaveRoot, srcFile))
	})

	t.Run("copy directory from a location to GameSaveRoot", func(t *testing.T) {
		rep := OSRepository{}
		srcDir := "test_dir"
		createDummyDirectory(t, srcDir)
		srcFile := path.Join(srcDir, "test_copy.txt")
		createDummyFile(t, srcFile)
		ensureCloned(t, emptyRepo)
		err := rep.Copy(srcDir, GameSaveRoot)
		assertNotError(t, err)
		assertExist(t, path.Join(GameSaveRoot, srcDir))
		assertExist(t, path.Join(GameSaveRoot, srcFile))
	})

	t.Run("copy with existing file", func(t *testing.T) {
		rep := OSRepository{}
		ensureCloned(t, emptyRepo)
		srcFile := "test_copy.txt"
		createDummyFile(t, srcFile)
		dstFile := path.Join(GameSaveRoot, srcFile)
		createBlankFile(t, dstFile)
		err := rep.Copy(srcFile, dstFile)
		assertNotError(t, err)
		assertSameContent(t, srcFile, dstFile)
	})

	t.Run("copy with undefined file", func(t *testing.T) {
		rep := OSRepository{}
		ensureCloned(t, emptyRepo)
		srcFile := "test_copy.txt"
		removeDummyFile(t, srcFile)
		err := rep.Copy(srcFile, GameSaveRoot)
		assertError(t, err)
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("Get existing config", func(t *testing.T) {
		rep := OSRepository{}
		initLocalConfig(t)
		addLocalConfig(t, "game_name", "game")
		gameName := rep.GetConfig("game_name")
		assertEqual(t, gameName, "game")
	})

	t.Run("Get undefined config", func(t *testing.T) {
		rep := OSRepository{}
		initLocalConfig(t)
		gameName := rep.GetConfig("game_name")
		assertEqual(t, gameName, "")
	})

	t.Run("Get config with undefined LocalConfig", func(t *testing.T) {
		rep := OSRepository{}
		removeLocalConfig(t)
		gameName := rep.GetConfig("game_name")
		assertEqual(t, gameName, "")
	})
}

func TestSetConfig(t *testing.T) {
	t.Run("Set existing config", func(t *testing.T) {
		rep := OSRepository{}
		initLocalConfig(t)
		addLocalConfig(t, "game_name", "game")
		err := rep.SetConfig("game_name", "new_game")
		assertNotError(t, err)
		val := getLocalConfig(t, "game_name")
		assertEqual(t, val, "new_game")
	})

	t.Run("Set undefined config", func(t *testing.T) {
		rep := OSRepository{}
		initLocalConfig(t)
		err := rep.SetConfig("game_name", "game")
		assertNotError(t, err)
		val := getLocalConfig(t, "game_name")
		assertEqual(t, val, "game")
	})

	t.Run("Set config with undefined LocalConfig", func(t *testing.T) {
		rep := OSRepository{}
		removeLocalConfig(t)
		err := rep.SetConfig("game_name", "game")
		assertNotError(t, err)
		assertExist(t, LocalConfig)
		val := getLocalConfig(t, "game_name")
		assertEqual(t, val, "game")
	})
}
