package command

import "errors"

var (
	errGameNotExist     = errors.New("Game is not exist, call add first")
	errGitInitialized   = errors.New("Git repo has been initialized")
	errGitUninitialized = errors.New("Git repo uninitialized, call init first")
	errSavePathNotExist = errors.New("Game save path is not exist, call set-path first")
)

type serviceMock struct {
	gameAdded    bool
	gamePrepared bool
	gitRepo      bool
	savePrepared bool
}

func newServiceMock() *serviceMock {
	return &serviceMock{
		gameAdded:    false,
		gamePrepared: false,
		gitRepo:      false,
		savePrepared: false,
	}
}

func (s *serviceMock) AddConfig(key, value string) error {
	if !s.gitRepo {
		return errGitUninitialized
	}
	if key == "save_path" {
		s.savePrepared = true
	} else if key == "game_name" {
		s.gameAdded = true
	}
	return nil
}

func (s *serviceMock) InitGitRepo(repoURL string) error {
	if s.gitRepo {
		return errGitInitialized
	}
	s.gitRepo = true
	return nil
}

func (s *serviceMock) LoadGame() error {
	if !s.gamePrepared {
		return errGameNotExist
	} else if !s.savePrepared {
		return errSavePathNotExist
	}
	return nil
}

func (s *serviceMock) PrepareGame() error {
	if !s.gameAdded {
		return errGameNotExist
	}
	if !s.gitRepo {
		return errGitUninitialized
	}
	s.gamePrepared = true
	return nil
}

func (s *serviceMock) SaveGame() error {
	if !s.gamePrepared {
		return errGameNotExist
	} else if !s.savePrepared {
		return errSavePathNotExist
	}
	return nil
}
