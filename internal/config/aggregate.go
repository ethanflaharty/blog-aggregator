package config

import (
	"context"
	"fmt"
)

func scrapeFeeds(s *State) error {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.DB.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	rssFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
