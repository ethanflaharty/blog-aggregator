package config

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/ethanflaharty/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	return nil
}
