package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/nico4565/blog-aggregator/internal/database"
)

func browseHandler(s *state, cmd command, user database.User) error {

	if len(cmd.args) > 1 {
		return fmt.Errorf("Error! Command usage: %s <limit>! Limit is an optional argument represents the max number of posts you want to fetch.", cmd.name)
	}

	limit := 2
	if len(cmd.args) != 0 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			log.Printf("Invalid limit value: %v", err)
			log.Print("Will use limit default")
		}
	}

	postSlice, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Coudn't fetch posts!\nErr: %v", err)
	}

	for i, post := range postSlice {
		fmt.Printf("-- post %d --\n", i)
		fmt.Println("")
		printPost(post)
		fmt.Println("")
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf("ID: %v\n", post.ID)
	fmt.Printf("CreatedAt: %v\n", post.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", post.UpdatedAt)
	fmt.Printf("Title: %v\n", post.Title)
	fmt.Printf("Url: %v\n", post.Url)
	fmt.Printf("Description: %v\n", post.Description)
	fmt.Printf("PublishedAt: %v\n", post.PublishedAt)
	fmt.Printf("FeedID: %v\n", post.FeedID)
}
