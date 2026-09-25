package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/silentfin/gator/internal/config"
	"github.com/silentfin/gator/internal/database"
	"github.com/silentfin/gator/internal/rss"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return val(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Empty args")
	}
	username := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		os.Exit(1)
		return err
	}

	if err := s.conf.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("User: %s has been set!\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	username := cmd.args[0]
	if len(cmd.args) == 0 {
		return fmt.Errorf("no name provided")
	}

	if _, err := s.db.GetUser(context.Background(), username); err == nil {
		os.Exit(1)
	}
	_, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username})
	if err != nil {
		return err
	}
	fmt.Printf("user created: %s\n", username)
	s.conf.SetUser(username)
	return nil
}

func handleReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
		return err
	}
	fmt.Println("reset successful!")
	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	curentUser := s.conf.CurrentUserName
	for _, user := range users {
		if user.Name == curentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handleAgg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("insufficient arguments, needs 2 args")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	fmt.Println("ID: ", feed.ID)
	fmt.Println("Name: ", feed.Name)
	fmt.Println("URL: ", feed.Url)
	fmt.Println("user_id: ", feed.UserID)
	fmt.Println("\nadded successful!")

	followCmd := command{
		name: "follow",
		args: []string{url},
	}
	err = handleFollow(s, followCmd, user)
	if err != nil {
		return err
	}
	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return nil
	}
	for _, feed := range feeds {
		name := feed.Name
		url := feed.Url
		user_id := feed.UserID
		userName, err := s.db.GetUserFromID(context.Background(), user_id)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s | URL: %s | userName: %s\n", name, url, userName.Name)

	}
	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("insufficient arguments, needs 1 args")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found or invalid url")
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("FeedName: %s\n", feedFollow.FeedName)
	fmt.Printf("UserName: %s\n", feedFollow.UserName)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s\n", feed.FeedName)
	}
	return nil
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	markedFetched, err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	feeds, err := rss.FetchFeed(context.Background(), markedFetched.Url)
	if err != nil {
		return err
	}

	var savedFeedCount int
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
	}

	for _, item := range feeds.Channel.Item {
		var pqErr *pq.Error
		var publishedAtTime time.Time
		var valid bool

		for _, layout := range layouts {
			timePublished, err := time.Parse(layout, item.PubDate)
			if err == nil {
				publishedAtTime = timePublished
				valid = true
				break
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: publishedAtTime, Valid: valid},
			FeedID:      markedFetched.ID,
		})
		if err != nil {
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			} else {
				log.Printf("error creating post: %v", err)
				continue
			}
		}
		savedFeedCount++
	}
	log.Printf("fetched feed %q saved %d new posts", markedFetched.Name, savedFeedCount)
	return nil
}

func handleBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
	}
	return nil
}

func handleHelp(s *state, cmd command) error {
	fmt.Println("Gator is a command-line RSS feed aggregator that fetches and stores posts.\n")
	fmt.Println(`Usage:
	gator <command> [arguments]
`)
	fmt.Println(`Available commands:
	register <name>          Create a user
	login <name>             Set the current user
	users                    List all users
	reset                    Delete all database data
	addfeed <name> <url>     Add and follow a feed
	feeds                    List feeds
	follow <url>             Follow an existing feed
	unfollow <url>           Unfollow a feed
	following                List followed feeds
	agg <duration>           Fetch feeds repeatedly, e.g. agg 10s
	browse [limit]           Show recent posts; defaults to 2
	help                     Show this message
`)
	return nil
}
