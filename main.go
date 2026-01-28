package main

import (
	"log"
	"os"

	"github.com/nico4565/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error! Couldn't read the config file: %v\n", err)
	}

	st := state{&cfg}

	cmds := commands{
		nameToHandlerMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)

	input := os.Args[:]
	if len(input) < 2 {
		log.Fatal("Not enough arguments, commands need a command name and 1 or more arguments!")
	}

	cmd := command{
		name: input[1],
		args: input[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatalf("Error:%s", err)
	}

}
