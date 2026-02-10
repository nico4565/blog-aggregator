package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nico4565/blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) > 1 {
		return fmt.Errorf("Login needs only one argument. Username field, no spaces!")
	}

	if len(cmd.args) < 1 {
		return fmt.Errorf("Login needs one argument. Username field, no spaces!")
	}

	cx := context.Background()

	_, err := s.db.GetUser(cx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("User not found! Register your user :)")
	}

	err = s.configPtr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set correctly!\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.args) > 1 {
		return fmt.Errorf("Register needs only one argument. Username field, no spaces!")
	}

	if len(cmd.args) < 1 {
		return fmt.Errorf("Register needs one argument. Username field, no spaces!")
	}

	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	name := cmd.args[0]

	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      name,
	}

	cx := context.Background()

	user, err := s.db.CreateUser(cx, params)
	if err != nil {
		return err
	}

	err = s.configPtr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been registered!\n", cmd.args[0])
	printUser(user)

	return nil
}

func handlerReset(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return fmt.Errorf("No arguments needed")
	}

	cx := context.Background()

	err := s.db.ResetUsers(cx)
	if err != nil {
		return fmt.Errorf("Error users not deleted!\nError:%s", err)
	}

	fmt.Println("All users and related user data have been deleted.")

	return nil
}

func handlerUsers(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return fmt.Errorf("No arguments needed")
	}

	cx := context.Background()

	users, err := s.db.GetUsers(cx)
	if err != nil {
		return fmt.Errorf("Error users not Found!\nError:%s", err)
	}

	if len(users) == 0 {
		fmt.Println("No users registered yet!")
		return nil
	}

	for _, user := range users {
		if user.Name == s.configPtr.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
