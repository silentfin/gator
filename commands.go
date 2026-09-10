package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/silentfin/gator/internal/config"
	"github.com/silentfin/gator/internal/database"
	"github.com/silentfin/gator/internal/rss"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return val(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Empty args")
	}
	username := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		os.Exit(1)
		return err
	}

	if err := s.conf.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User: %s has been set!\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	username := cmd.args[0]
	if len(cmd.args) == 0 {
		return fmt.Errorf("no name provided")
	}

	if _, err := s.db.GetUser(context.Background(), username); err == nil {
		os.Exit(1)
	}
	_, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username})
	if err != nil {
		return err
	}
	fmt.Printf("user created: %s\n", username)
	s.conf.SetUser(username)
	return nil
}

func handleReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		os.Exit(1)
		return err
	}
	fmt.Println("reset successful!")
	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	curentUser := s.conf.CurrentUserName
	for _, user := range users {
		if user.Name == curentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	url := cmd.args[0]
	feedData, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feedData)
	return nil
}
