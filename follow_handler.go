package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nico4565/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("Follow needs only 1 argument!(the url of the RSS you want to follow)")
	}

	if len(cmd.args) < 1 {
		return fmt.Errorf("Follow needs only 1 argument!(the url of the RSS you want to follow)")
	}

	cx := context.Background()

	feed, err := s.db.GetFeedByUrl(cx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("Feed not found, url missing. Maybe you should add this feed yourself!")
	}

	user, err := s.db.GetUser(cx, s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("User not found!")
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(cx, params)
	if err != nil {
		return fmt.Errorf("Something went wrong!\nErr:%w\n", err)
	}

	fmt.Println("FeedFollow created successfully:")
	fmt.Println("")
	printFeedFollowRow(feedFollow)
	fmt.Println()
	fmt.Println("")

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("Following doesn't need args!")
	}

	cx := context.Background()

	user, err := s.db.GetUser(cx, s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("User not found!")
	}

	currentUserFeeds, err := s.db.GetFeedFollowByUser(cx, user.ID)
	if err != nil {
		return err
	}

	if len(currentUserFeeds) == 0 {
		fmt.Println("Current user follows 0 feeds.")
		return nil
	}

	for i, userFeed := range currentUserFeeds {
		fmt.Printf("-- feed %d --\n", i)
		fmt.Println("")
		printFeedFollowRowByUser(userFeed)
		fmt.Println()
		fmt.Println("")
	}

	return nil
}

func printFeedFollowRow(row database.CreateFeedFollowRow) {
	fmt.Printf("* ID:            %s\n", row.ID)
	fmt.Printf("* Created:       %v\n", row.CreatedAt)
	fmt.Printf("* Updated:       %v\n", row.UpdatedAt)
	fmt.Printf("* UserID:        %s\n", row.UserID)
	fmt.Printf("* FeedID:        %s\n", row.FeedID)
	fmt.Printf("* User:        %s\n", row.User)
	fmt.Printf("* Feed:        %s\n", row.Feed)
}

func printFeedFollowRowByUser(row database.GetFeedFollowByUserRow) {
	fmt.Printf("* ID:            %s\n", row.ID)
	fmt.Printf("* Created:       %v\n", row.CreatedAt)
	fmt.Printf("* Updated:       %v\n", row.UpdatedAt)
	fmt.Printf("* UserID:        %s\n", row.UserID)
	fmt.Printf("* FeedID:        %s\n", row.FeedID)
	fmt.Printf("* User:        %s\n", row.User)
	fmt.Printf("* Feed:        %s\n", row.Feed)
}
