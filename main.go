package main

import (
	"fmt"
	"os"

	"github.com/ethanflaharty/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var cfgFile config.State
	cfgFile.CFG = &cfg

	var savedCmds config.Commands
	savedCmds.Cmds = make(map[string]func(*config.State, config.Command) error)
	savedCmds.Register("login", config.HandlerLogin)
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	cmdName := os.Args[1]
	cmdUsername := os.Args[2:]
	var newCmd config.Command
	newCmd.Name = cmdName
	newCmd.Username = cmdUsername

	err = savedCmds.Run(&cfgFile, newCmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
