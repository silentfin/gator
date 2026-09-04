package main

import (
	"fmt"

	"github.com/silentfin/gator/internal/config"
)

type state struct {
	conf *config.Config
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
	if err := s.conf.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User: %s has been set!\n", cmd.args[0])
	return nil
}
