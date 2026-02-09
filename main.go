package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nico4565/blog-aggregator/internal/config"
	"github.com/nico4565/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error! Couldn't read the config file: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)

	st := state{dbQueries, &cfg}

	cmds := commands{
		nameToHandlerMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)

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
