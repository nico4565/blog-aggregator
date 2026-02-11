package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nico4565/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return fmt.Errorf("Error! Command usage: %s <time_duration>! time_duration examples: \"1m\", \"1m1s\",\"1h\" and so on.", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s", cmd.args[0])

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	cx := context.Background()
	feedToFetch, err := s.db.GetNextFeedToFetch(cx)
	if err != nil {
		return fmt.Errorf("Error: Couldn't get next feed to fetch!")
	}

	feed_id := feedToFetch.ID

	err = s.db.MarkFeedFetched(cx, feed_id)
	if err != nil {
		return fmt.Errorf("Error: Couldn't get mark feed %s with id \"%s\"as fetched!", feedToFetch.Name, feedToFetch.ID)
	}

	feedData, err := fetchFeed(cx, feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed!\n%s", err)
	}

	postsStoredCounter := 0

	for _, item := range feedData.Channel.Item {

		fmt.Printf("Found post: %s\n", item.Title)
		t := parseStringToTime(item.PubDate)
		_, err = s.db.StorePost(cx, database.StorePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: item.Title != ""},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: t, Valid: t != time.Time{}},
			FeedID:      feed_id,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
				postsStoredCounter++
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
		postsStoredCounter++
	}
	log.Printf("Feed %s collected\nPosts found: %d\nPosts stored: %d\n", feedToFetch.Name, len(feedData.Channel.Item), postsStoredCounter)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.args) != 2 {
		return fmt.Errorf("Error! Command usage: %s <feed_name> <feed_url>", cmd.name)
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

	params := database.StoreFeedParams{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feedSt, err := s.db.StoreFeed(cx, params)
	if err != nil {
		return err
	}

	fmt.Println("Feed stored successfully.")

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedSt.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("FeedFollow created successfully:")
	fmt.Println("")
	printFeedFollowRow(feedFollow)
	fmt.Println()
	fmt.Println("")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {

	if len(cmd.args) > 0 {
		return fmt.Errorf("No arguments needed")
	}

	cx := context.Background()

	feeds, err := s.db.GetFeeds(cx)
	if err != nil {
		return fmt.Errorf("Error feeds not Found!\nError:%s", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds stored yet!")
		return nil
	}

	for i, feed := range feeds {
		fmt.Printf("-- feed %d --\n", i)
		fmt.Println("")

		printFeedEntity(feed)
		user, err := s.db.GetUserById(cx, feed.UserID)
		if err != nil {
			return fmt.Errorf("GetUserById failed!\nError:%s", err)
		}
		fmt.Printf("UserName: %v\n", user.Name)
		fmt.Println("")
	}

	return nil
}

func parseStringToTime(timeString string) time.Time {
	layout1 := time.RFC1123Z
	layout2 := time.RFC1123

	t, err := time.Parse(layout1, timeString)
	if err != nil {
		t, err = time.Parse(layout2, timeString)
		if err != nil {
			log.Printf("couldn't parse date %q: %v", timeString, err)
			return time.Time{}
		}
		return t
	}

	return t
}

func printFeedEntity(feed database.Feed) {
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserId: %v\n", feed.UserID)
	fmt.Printf("LastFetchedAt: %v\n", feed.LastFetchedAt)
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
