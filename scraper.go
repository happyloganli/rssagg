package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/happyloganli/rssagg/internal/database"
	"log"
	"strings"
	"sync"
	"time"
)

func startScraping(
	queries *database.Queries,
	concurrency int,
	timeBetweenScraping time.Duration,
) {

	log.Printf("Starting scraping with concurrency %d every %s duration", concurrency, timeBetweenScraping)
	ticker := time.NewTicker(timeBetweenScraping)
	for ; ; <-ticker.C {
		feeds, err := queries.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, &feed, queries)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, feed *database.Feed, queries *database.Queries) {
	defer wg.Done()

	_, err := queries.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Error parsing feed url: %v", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing feed pub date: %v", err)
		}

		_, err = queries.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Description: description,
				PublishedAt: pubAt,
				Url:         item.Link,
				FeedID:      feed.ID,
			})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post: %v", err)
		}
	}

	log.Printf("Feeds %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
