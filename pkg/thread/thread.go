package thread

import (
	"fmt"
	"strings"

	"github.com/dkaslovsky/thread-safe/pkg/twitter"
)

const (
	// maxThreadLen is the maximum number of tweets to be fetched for constructing a thread
	maxThreadLen = 100
)

// Thread represents a Twitter thread
type Thread struct {
	Name   string           `json:"name"`
	Tweets []*twitter.Tweet `json:"tweets"`
}

// New constructs a Thread by querying the Twitter API
func New(client twitter.Client, name string, lastID string) (*Thread, error) {
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

// Len returns the number of tweets contained in a Thread
func (th *Thread) Len() int {
	return len(th.Tweets)
}

// Metadata returns a string with thread metadata
func (th *Thread) Metadata() string {
	if th.Len() == 0 {
		return ""
	}
	first := th.Tweets[0]
	headerStrs := []string{
		fmt.Sprintf("URL: \t\t\t%s", first.URL),
		fmt.Sprintf("Author Name: \t\t%s", first.AuthorName),
		fmt.Sprintf("Author Handle: \t\t%s", first.AuthorHandle),
		fmt.Sprintf("Conversation ID: \t%s", first.ConversationID),
	}
	return strings.Join(headerStrs, "\n")
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

func reverseSlice[T any](s []T) {
	first, last := 0, len(s)-1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}
