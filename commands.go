package main

import (
	"fmt"

	"github.com/nico4565/blog-aggregator/internal/config"
	"github.com/nico4565/blog-aggregator/internal/database"
)

type state struct {
	db        *database.Queries
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

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.nameToHandlerMap[name] = f
}
