package main

import (
	"context"
	"database/sql"
	"fmt"
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
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.configPtr.CurrentUserName)
		if err != nil {
			return fmt.Errorf("User not found!")
		}

		return handler(s, cmd, user)
	}
}
