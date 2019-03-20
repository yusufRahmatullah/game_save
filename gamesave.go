package main

import "github.com/yusufRahmatullah/game_save/command"

func main() {
	root := command.NewRootCommand(nil, nil)
	root.Run()
}
