package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/g8rswimmer/go-twitter/v2"
)

// maxThreadLen is the maximum number of tweets allowed to be fetched for constructing a thread
const maxThreadLen = 100

type threadSaver struct {
	client Client
	opts   twitter.TweetLookupOpts
}

func newThreadSaver(token string, opts twitter.TweetLookupOpts) *threadSaver {
	return &threadSaver{
		client: newClient(token),
		opts:   opts,
	}
}

func (ts *threadSaver) thread(lastID string) (*thread, error) {
	tweets, err := ts.walkTweets(lastID, maxThreadLen)
	if err != nil {
		return nil, err
	}

	// Tweets are fetched from last to first so reverse the order
	reverseSlice(tweets)

	return &thread{tweets}, nil
}

func (ts *threadSaver) tweet(id string) (*tweet, error) {
	return ts.client.tweetLookup(id, ts.opts)
}

func (ts *threadSaver) walkTweets(id string, limit int) ([]*tweet, error) {
	tweets := []*tweet{}

	nextID := id
	conversationID := ""
	authorID := ""

	for i := 0; i < limit; i++ {
		tweet, err := ts.tweet(nextID)
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
		// Next ID for lookup
		case 1:
			nextID = tweet.RepliedToIDs[0]
		// Top of thread has been reached
		case 0:
			return tweets, nil
		// Error on multiple replied_to IDs
		default:
			return nil, fmt.Errorf("cannot follow tweet %s with multiple replied_to IDs", tweet.ID)
		}
	}

	log.Printf("exceeded maximum [%d] number of tweets to fetch", limit)
	return tweets, nil
}

func reverseSlice[T any](s []T) {
	first, last := 0, len(s)-1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}

func main() {
	token := flag.String("token", "", "twitter API bearer token") // TODO: read from env
	id := flag.String("id", "", "id of the last tweet in a single-author thread")
	flag.Parse()

	ts := newThreadSaver(*token, twitter.TweetLookupOpts{
		Expansions: []twitter.Expansion{
			twitter.ExpansionEntitiesMentionsUserName,
			twitter.ExpansionAuthorID,
			twitter.ExpansionAttachmentsMediaKeys,
		},
		MediaFields: []twitter.MediaField{
			twitter.MediaFieldMediaKey,
			twitter.MediaFieldURL,
			twitter.MediaFieldType,
			twitter.MediaFieldPreviewImageURL,
		},
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldReferencedTweets,
		},
	})

	thread, err := ts.thread(*id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(thread)
}
