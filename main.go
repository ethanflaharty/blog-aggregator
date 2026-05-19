package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ethanflaharty/blog-aggregator/internal/config"
	"github.com/ethanflaharty/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dbQueries := database.New(db)
	var cfgFile config.State
	cfgFile.CFG = &cfg
	cfgFile.DB = dbQueries

	var savedCmds config.Commands
	savedCmds.Cmds = make(map[string]func(*config.State, config.Command) error)
	savedCmds.Register("login", config.HandlerLogin)
	savedCmds.Register("register", config.HandlerRegister)
	savedCmds.Register("reset", config.HandlerReset)
	savedCmds.Register("users", config.HandlerUsers)
	savedCmds.Register("agg", config.HandlerAgg)
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
