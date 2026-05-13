package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethanflaharty/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	DB  *database.Queries
}

type Command struct {
	Name     string
	Username []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Username) < 1 {
		return errors.New("login expects a username")
	}
	_, err := s.DB.GetUser(context.Background(), cmd.Username[0])
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
	err = s.CFG.SetUser(cmd.Username[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Username) < 1 {
		return errors.New("register expects a username")
	}
	_, err := s.DB.GetUser(context.Background(), cmd.Username[0])
	if err == nil {
		return fmt.Errorf("user already exists")
	}
	userParms := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Username[0]}
	newUser, err := s.DB.CreateUser(context.Background(), userParms)
	if err != nil {
		return err
	}
	err = s.CFG.SetUser(newUser.Name)
	if err != nil {
		return err
	}
	fmt.Println("new user was created")
	fmt.Println(newUser)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DB.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("reset unsuccessful")
	}
	fmt.Println("reset successful")
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
