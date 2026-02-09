package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/nico4565/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return fmt.Errorf("No arguments needed")
	}

	ctx := context.Background()

	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error fetching feed!\n%s", err)
	}

	printRSSFeed(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.args) > 2 {
		return fmt.Errorf("AddFeed needs only 2 arguments. Rss name and url!")
	}

	if len(cmd.args) < 2 {
		return fmt.Errorf("AddFeed needs 2 argument. Rss name and url!")
	}

	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	name := cmd.args[0]
	u, err := url.Parse(cmd.args[1])
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return err
	}
	url := cmd.args[1]

	cx := context.Background()
	user, err := s.db.GetUser(cx, s.configPtr.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error users not Found!\nError:%s", err)
	}

	params := database.StoreFeedParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      name,
		Url:       url,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	}

	feed, err := s.db.StoreFeed(cx, params)
	if err != nil {
		return err
	}

	printFeedEntity(feed)

	return nil
}

func printFeedEntity(feed database.Feed) {
	fmt.Printf("Printing feed stored")
	fmt.Printf("ID: %v", feed.ID)
	fmt.Printf("CreatedAt: %v", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v", feed.UpdatedAt)
	fmt.Printf("Name: %v", feed.Name)
	fmt.Printf("Url: %v", feed.Url)
	fmt.Printf("UserId: %v", feed.UserID)
}

func printRSSFeed(feed *RSSFeed) {
	fmt.Printf(" * Title:      %v\n", feed.Channel.Title)
	fmt.Printf(" * Link:      %v\n", feed.Channel.Link)
	fmt.Printf(" * Description:    %v\n", feed.Channel.Description)
	for i := range feed.Channel.Item {
		fmt.Printf("Item %v:\n", i)
		printItem(feed.Channel.Item[i])
	}
}

func printItem(item RSSItem) {
	fmt.Printf(" 	* ItemTitle:      %v\n", item.Title)
	fmt.Printf(" 	* ItemLink:      %v\n", item.Link)
	fmt.Printf(" 	* itemDescription:    %v\n", item.Description)
	fmt.Printf(" 	* ItemPubDate:      %v\n", item.PubDate)
}
