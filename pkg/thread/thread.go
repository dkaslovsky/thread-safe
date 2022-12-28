package thread

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dkaslovsky/thread-safe/pkg/twitter"
)

const (
	// maxThreadLen is the maximum number of tweets allowed to be fetched for constructing a thread
	maxThreadLen = 100
)

type Thread struct {
	tweets []*twitter.Tweet
}

func NewThread(client twitter.Client, lastID string) (*Thread, error) {
	tweets, err := walkTweets(client, lastID, maxThreadLen)
	if err != nil {
		return nil, err
	}

	// Tweets are fetched from last to first so reverse the order
	reverseSlice(tweets)

	return &Thread{tweets}, nil
}

func NewThreadFromFile(path string) (*Thread, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := &Thread{}
	return t, json.Unmarshal(b, &(t.tweets))
}

func (t *Thread) ToFile(path string) error {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return err
	}
	b, berr := json.Marshal(t.tweets)
	if berr != nil {
		return berr
	}
	return os.WriteFile(filepath.Join(path, "thread.json"), b, 0o755)
}

func (t *Thread) DownloadAttachments(path string) error {
	for _, tweet := range t.Tweets() {
		for _, attachment := range tweet.Attachments {
			err := attachment.Download(filepath.Join(path, attachment.Name(tweet.ID)))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Thread) Tweets() []*twitter.Tweet {
	return t.tweets
}

func (t *Thread) Len() int {
	return len(t.tweets)
}

func (t *Thread) Header() string {
	if t.Len() == 0 {
		return ""
	}
	first := t.Tweets()[0]
	headerStrs := []string{
		fmt.Sprintf("URL: \t\t\t%s", first.URL), // TODO: sanitize HTML here or in template
		fmt.Sprintf("Author Name: \t\t%s", first.AuthorName),
		fmt.Sprintf("Author Handle: \t\t%s", first.AuthorHandle),
		fmt.Sprintf("Conversation ID: \t%s", first.ConversationID),
	}
	return strings.Join(headerStrs, "\n")
}

func walkTweets(client twitter.Client, id string, limit int) ([]*twitter.Tweet, error) {
	tweets := []*twitter.Tweet{}

	nextID := id
	conversationID := ""
	authorID := ""

	for i := 0; i < limit; i++ {
		tweet, err := client.LookupTweet(nextID)
		if err != nil {
			return nil, err
		}

		// Save the conversationID and authorID
		if i == 0 {
			conversationID = tweet.ConversationID
			authorID = tweet.AuthorID
		}

		// A change in conversationID or authorID indicates the end of the current thread
		if tweet.ConversationID != conversationID || tweet.AuthorID != authorID {
			return tweets, nil
		}

		tweets = append(tweets, tweet)

		switch len(tweet.RepliedToIDs) {
		case 1: // Next ID for lookup
			nextID = tweet.RepliedToIDs[0]
		case 0: // Top of thread has been reached
			return tweets, nil
		default: // Error on multiple replied_to IDs
			return nil, fmt.Errorf("cannot follow tweet %s with multiple replied_to IDs", tweet.ID)
		}
	}

	// Limit reached
	return nil, fmt.Errorf("exceeded maximum number of tweets to fetch [%d]", limit)
}

func reverseSlice[T any](s []T) {
	first, last := 0, len(s)-1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}
