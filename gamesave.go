package main

import (
	"github.com/yusufRahmatullah/game_save/command"
	"github.com/yusufRahmatullah/game_save/repository"
	"github.com/yusufRahmatullah/game_save/service"
)

const (
	// AppVersion is the version of GameSave
	AppVersion = "0.1.0"
)

func main() {
	gitRepo := repository.GitRepository{}
	osRepo := repository.OSRepository{}
	service := service.Service{
		GitRepository: &gitRepo,
		OSRepository:  &osRepo,
	}
	root := command.NewRootCommand(&service)
	root.SetVersion(AppVersion)
	root.Run()
}
