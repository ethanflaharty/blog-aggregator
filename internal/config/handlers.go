package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethanflaharty/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.CFG.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Username) < 2 {
		return fmt.Errorf("not enough args provided")
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Username[0],
		Url:       cmd.Username[1],
		UserID:    user.ID,
	}
	feed, err := s.DB.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feedFollow.FeedName)
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

func HanderFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Username) < 1 {
		return fmt.Errorf("no url provided")
	}
	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Username[0])
	if err != nil {
		return err
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	row, err := s.DB.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Println(row.FeedName)
	fmt.Println(user.Name)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	following, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range following {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Username) < 1 {
		return fmt.Errorf("requires feed url")
	}
	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Username[0])
	if err != nil {
		return err
	}
	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.DB.UnfollowFeed(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}
