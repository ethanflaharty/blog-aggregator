package config

import (
	"errors"
	"fmt"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

type State struct {
	CFG *Config
}

type Command struct {
	Name     string
	Username []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Username) < 1 {
		return errors.New("login expects a username")
	}
	err := s.CFG.SetUser(cmd.Username[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Cmds[cmd.Name]
	if !ok {
		return errors.New("command doesn't exist")
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}
