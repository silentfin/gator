package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/silentfin/gator/internal/config"
	"github.com/silentfin/gator/internal/database"
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
