package service

import (
	"errors"
	"fmt"

	"github.com/yusufRahmatullah/game_save/repository"
)

var (
	// ErrGameNameEmpty represents error if Game name has not been set
	ErrGameNameEmpty = errors.New("Game name has not been set")
	// ErrSavePathEmpty represents error if Game save path has not been set
	ErrSavePathEmpty = errors.New("Game save path has not been set")
)

// IService is interface for interaction with repositories
type IService interface {
	AddConfig(key, value string) error
	InitGitRepo(repoURL string) error
	LoadGame() error
	LoadGameSave() error
	SaveGame() error
}

// Service is the implementation of IService
type Service struct {
	GitRepository repository.IGitRepository
	OSRepository  repository.IOSRepository
}

// AddConfig add key and value to configuration
func (s *Service) AddConfig(key, value string) error {
	return s.OSRepository.SetConfig(key, value)
}

// InitGitRepo initialize Git repository URL
func (s *Service) InitGitRepo(repoURL string) error {
	return s.GitRepository.Clone(repoURL)
}

// LoadGame prepare Git to change the current branch to game name
func (s *Service) LoadGame() error {
	gameName := s.OSRepository.GetConfig("game_name")
	if gameName == "" {
		return ErrGameNameEmpty
	}
	err := s.GitRepository.Checkout(gameName)
	if err != nil {
		return err
	}
	err = s.GitRepository.Pull(gameName)
	if err != nil {
		return err
	}
	return nil
}

// LoadGameSave load game's save data by copying the save data
// from git repository to save path
func (s *Service) LoadGameSave() error {
	savePath := s.OSRepository.GetConfig("save_path")
	if savePath == "" {
		return ErrSavePathEmpty
	}
	return s.OSRepository.Copy(repository.GameSaveRoot, savePath)
}

// SaveGame persists game's save data by copying save data
// from save path to git repository
func (s *Service) SaveGame() error {
	gameName := s.OSRepository.GetConfig("game_name")
	if gameName == "" {
		return ErrGameNameEmpty
	}
	savePath := s.OSRepository.GetConfig("save_path")
	if savePath == "" {
		return ErrSavePathEmpty
	}
	err := s.OSRepository.Copy(savePath, repository.GameSaveRoot)
	if err != nil {
		return err
	}
	err = s.GitRepository.Checkout(gameName)
	if err != nil {
		return err
	}
	err = s.GitRepository.Commit(s.generateCommitMessage())
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) generateCommitMessage() string {
	gameName := s.OSRepository.GetConfig("game_name")
	return fmt.Sprintf("Update %s", gameName)
}
