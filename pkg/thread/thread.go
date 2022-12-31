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
	// attachmentsDirName is the name used for the directory where attachment files are saved
	attachmentsDirName = "attachments"
	// jsonFileName is the name used for the generated JSON file
	jsonFileName = "thread.json"
)

// Thread represents a Twitter thread
type Thread struct {
	Name   string           `json:"name"`
	Tweets []*twitter.Tweet `json:"tweets"`
}

// NewThread constructs a Thread by querying the Twitter API
func NewThread(client twitter.Client, name string, lastID string) (*Thread, error) {
	tweets, err := walkTweets(client, lastID, maxThreadLen)
	if err != nil {
		return nil, err
	}

	// Tweets are fetched from last to first so reverse the order
	reverseSlice(tweets)

	return &Thread{
		Name:   name,
		Tweets: tweets,
	}, nil
}

// NewThreadFromFile constructs a Thread by loading data from a file
func NewThreadFromFile(path string) (*Thread, error) {
	filePath := filepath.Clean(filepath.Join(path, jsonFileName))
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	th := &Thread{}
	return th, json.Unmarshal(b, &th)
}

// ToJSON generates and saves a JSON file from a Thread's tweets
func (th *Thread) ToJSON(path string) error {
	err := os.MkdirAll(filepath.Clean(path), 0o750)
	if err != nil {
		return err
	}

	b, bErr := json.Marshal(th)
	if bErr != nil {
		return bErr
	}

	jsonPath := filepath.Clean(filepath.Join(path, jsonFileName))
	return os.WriteFile(jsonPath, b, 0o600)
}

// DownloadAttachments saves all media attachments from a Thread's
func (th *Thread) DownloadAttachments(path string) error {
	attachmentPath := filepath.Join(path, attachmentsDirName)
	err := os.MkdirAll(attachmentPath, 0o750)
	if err != nil {
		return err
	}

	for _, tweet := range th.Tweets {
		for _, attachment := range tweet.Attachments {
			attachmentName := attachment.Name(tweet.ID)
			err := attachment.Download(filepath.Clean(filepath.Join(attachmentPath, attachmentName)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Len returns the number of tweets contained in a Thread
func (th *Thread) Len() int {
	return len(th.Tweets)
}

// walkTweets queries for tweets by following the RepliedToID of the starting tweet and stopping
// once no more tweets are in the chain or a new conversation ID or author ID is encountered
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

// Dir constructs the name of the directory where thread files are written
func Dir(topLevelPath string, threadName string) string {
	return filepath.Join(topLevelPath, strings.Replace(threadName, " ", "_", -1))
}

func reverseSlice[T any](s []T) {
	first, last := 0, len(s)-1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}
