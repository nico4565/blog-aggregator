package main

import (
	"fmt"

	"github.com/nico4565/blog-aggregator/internal/config"
)

type state struct {
	configPtr *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	nameToHandlerMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.nameToHandlerMap[cmd.name]
	if !exists {
		return fmt.Errorf("Command %s doesn't exist :(", cmd.name)
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.nameToHandlerMap[name] = f
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) > 1 {
		return fmt.Errorf("Login needs only one string argument. Username field, no spaces!")
	}

	if len(cmd.args) < 1 {
		return fmt.Errorf("Login needs one string argument. Username field, no spaces!")
	}

	err := s.configPtr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set!\n", cmd.args[0])

	return nil
}
