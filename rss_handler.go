package main

import (
	"context"
	"fmt"
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

	printFeed(feed)

	return nil
}

func printFeed(feed *RSSFeed) {
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
