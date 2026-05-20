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

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("no users saved")
	}
	for _, name := range users {
		if name != s.CFG.CurrentUserName {
			fmt.Printf("* %v\n", name)
		} else {
			fmt.Printf("* %v (current)\n", name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Username) < 2 {
		return fmt.Errorf("not enough args provided")
	}
	user, err := s.DB.GetUser(context.Background(), s.CFG.CurrentUserName)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   cmd.Username[0],
		Url:    cmd.Username[1],
		UserID: user.ID,
	}
	feed, err := s.DB.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2)
	}
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
