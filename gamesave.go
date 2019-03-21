package main

import (
	"github.com/yusufRahmatullah/game_save/command"
	"github.com/yusufRahmatullah/game_save/repository"
	"github.com/yusufRahmatullah/game_save/service"
)

func main() {
	gitRepo := repository.GitRepository{}
	osRepo := repository.OSRepository{}
	service := service.Service{
		GitRepository: &gitRepo,
		OSRepository:  &osRepo,
	}
	root := command.NewRootCommand(&service)
	root.Run()
}
