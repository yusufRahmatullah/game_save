package service

import "testing"

func TestAddConfig(t *testing.T) {
	t.Run("set game_name configuration", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		err := service.AddConfig("game_name", "game")
		assertNotError(t, err)
		gameName := service.OSRepository.GetConfig("game_name")
		assertEqual(t, gameName, "game")
	})
}

func TestInitGitRepo(t *testing.T) {
	t.Run("initialize git repository using valid URL", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		err := service.InitGitRepo("dummy.git")
		assertNotError(t, err)
	})

	t.Run("initialize git repository using invalid URL", func(t *testing.T) {
		service := initService(t, gitOptionRepoInvalid)
		err := service.InitGitRepo("dummy.git")
		assertError(t, err)
	})
}

func TestLoadGame(t *testing.T) {
	t.Run("load game in normal condition", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		service.AddConfig("game_name", "game")
		err := service.LoadGame()
		assertNotError(t, err)
		currentBranch, _ := service.GitRepository.GetCurrentBranch()
		assertEqual(t, currentBranch, "game")
	})

	t.Run("load game branch with game name not exist", func(t *testing.T) {
		service := initService(t, gitOptionsBranchInvalid)
		service.AddConfig("game_name", "game")
		err := service.LoadGame()
		assertNotError(t, err)
		currentBranch, _ := service.GitRepository.GetCurrentBranch()
		assertEqual(t, currentBranch, "game")
	})

	t.Run("load game game_name not set", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		err := service.LoadGame()
		assertError(t, err)
	})
}

func TestLoadGameSave(t *testing.T) {
	t.Run("load game save in normal condition", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		service.AddConfig("save_path", "./game.save")
		err := service.LoadGameSave()
		assertNotError(t, err)
	})

	t.Run("save_path not set", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		err := service.LoadGameSave()
		assertError(t, err)
	})
}

func TestSaveGame(t *testing.T) {
	t.Run("save game in normal condition", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		service.AddConfig("save_path", "./game.save")
		service.AddConfig("game_name", "game")
		err := service.SaveGame()
		assertNotError(t, err)
	})

	t.Run("game name not set", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		service.AddConfig("save_path", "./game.save")
		err := service.SaveGame()
		assertError(t, err)
	})

	t.Run("save path not set", func(t *testing.T) {
		service := initService(t, gitOptionNormal)
		service.AddConfig("game_name", "game")
		err := service.SaveGame()
		assertError(t, err)
	})

	t.Run("git repo url not set", func(t *testing.T) {
		service := initService(t, gitOptionRepoInvalid)
		service.AddConfig("save_path", "./game.save")
		service.AddConfig("game_name", "game")
		err := service.SaveGame()
		assertError(t, err)
	})
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("Got '%s' expect '%s'", got, want)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Should be thrown an error")
	}
}

func assertNotError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Should be not error. Error: %v", err)
	}
}

func initService(t *testing.T, gitOptions map[string]bool) *Service {
	t.Helper()
	return &Service{
		GitRepository: NewGitRepositoryMock(gitOptions),
		OSRepository:  NewOsRepositoryMock(),
	}
}
