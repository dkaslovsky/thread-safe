package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
)

const MaxThreadLen = 100

type ThreadSaver struct {
	client *twitter.Client
	opts   twitter.TweetLookupOpts
}

func NewThreadSaver(token string, opts twitter.TweetLookupOpts) *ThreadSaver {
	return &ThreadSaver{
		client: &twitter.Client{
			Authorizer: authorize{
				Token: token,
			},
			Client: http.DefaultClient,
			Host:   "https://api.twitter.com",
		},
		opts: opts,
	}
}

func (ts *ThreadSaver) GetThread(lastID string) (Thread, error) {
	prevID := lastID
	conversationID := ""
	authorID := ""

	tweets := []Tweet{}
	for i := 0; i < MaxThreadLen; i++ {
		tweet, err := ts.getTweet(prevID)
		if err != nil {
			return Thread{}, err
		}
		tweets = append(tweets, tweet)

		if tweet.RepliedToID == "" {
			break
		}

		if conversationID == "" {
			conversationID = tweet.ConversationID
		}
		if authorID == "" {
			authorID = tweet.AuthorID
		}
		if tweet.ConversationID != conversationID || tweet.AuthorID != authorID {
			// TODO: log...
			break
		}

		prevID = tweet.RepliedToID
	}

	// Reverse the order
	i, j := 0, len(tweets)-1
	for i < j {
		tweets[i], tweets[j] = tweets[j], tweets[i]
		i++
		j--
	}

	return Thread{
		Tweets: tweets,
	}, nil
}

func (ts *ThreadSaver) getTweet(tweetID string) (Tweet, error) {
	tweetResponse, err := ts.client.TweetLookup(context.Background(), []string{tweetID}, ts.opts)
	if err != nil {
		return Tweet{}, fmt.Errorf("tweet lookup error: %v", err)
	}

	tweetDictionary, ok := tweetResponse.Raw.TweetDictionaries()[tweetID]
	if !ok {
		return Tweet{}, fmt.Errorf("tweet lookup error: response does not include tweet with ID %s", tweetID)
	}

	return ts.parseTweet(tweetDictionary)
}

func (ts *ThreadSaver) parseTweet(raw *twitter.TweetDictionary) (Tweet, error) {
	repliedToID := ""
	for _, ref := range raw.Tweet.ReferencedTweets {
		if ref.Type == "replied_to" {
			repliedToID = ref.ID
			break
		}
	}

	return Tweet{
		ID:             raw.Tweet.ID,
		ConversationID: raw.Tweet.ConversationID,
		URL:            "", // TODO
		Text:           raw.Tweet.Text,
		CreatedAt:      raw.Tweet.CreatedAt,
		AuthorID:       raw.Tweet.AuthorID,
		AuthorName:     raw.Author.Name,
		AuthorHandle:   raw.Author.UserName,
		RepliedToID:    repliedToID,
	}, nil
}

type Thread struct {
	Tweets []Tweet
}

func (t Thread) Len() int {
	return len(t.Tweets)
}

func (t Thread) String() string {
	tweetStrs := []string{}
	threadLen := t.Len()
	for i, tweet := range t.Tweets {
		tweetStr := fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text)
		tweetStrs = append(tweetStrs, tweetStr)
	}
	return strings.Join(tweetStrs, "\n---\n")
}

type Tweet struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	URL            string `json:"url"`
	Text           string `json:"text"`
	CreatedAt      string `json:"created_at"`
	AuthorID       string `json:"author_id"`
	AuthorName     string `json:"author_name"`
	AuthorHandle   string `json:"author_handle"`
	RepliedToID    string `json:"replied_to_id"`
	// Attachments  []Attachment // TODO
}

// func (t *Tweet) String() string {
// 	b, _ := json.MarshalIndent(t, "", "\t")
// 	return string(b)
// }

// type Attachment struct {
// 	MediaKey string
// 	Type     string
// 	URL      string
// }

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func main() {
	token := flag.String("token", "", "twitter API bearer token")
	id := flag.String("id", "", "id of the last tweet in a single-author thread")
	flag.Parse()

	ts := NewThreadSaver(*token, twitter.TweetLookupOpts{
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

	thread, err := ts.GetThread(*id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(thread)
}
