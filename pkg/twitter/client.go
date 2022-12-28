package twitter

import (
	"context"
	"fmt"
	"net/http"

	tw "github.com/g8rswimmer/go-twitter/v2"
)

type Client interface {
	LookupTweet(id string) (*Tweet, error)
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func NewTwitterClient(token string) Client {
	return &client{
		&tw.Client{
			Authorizer: authorize{
				Token: token,
			},
			Client: http.DefaultClient,
			Host:   "https://api.twitter.com",
		},
	}
}

type client struct {
	*tw.Client
}

func (c *client) LookupTweet(tweetID string) (*Tweet, error) {
	tweetResponse, err := c.TweetLookup(context.Background(), []string{tweetID}, tw.TweetLookupOpts{
		Expansions: []tw.Expansion{
			tw.ExpansionEntitiesMentionsUserName,
			tw.ExpansionAuthorID,
			tw.ExpansionAttachmentsMediaKeys,
		},
		MediaFields: []tw.MediaField{
			tw.MediaFieldMediaKey,
			tw.MediaFieldURL,
			tw.MediaFieldType,
			tw.MediaFieldPreviewImageURL,
			tw.MediaFieldVariants,
		},
		TweetFields: []tw.TweetField{
			tw.TweetFieldCreatedAt,
			tw.TweetFieldConversationID,
			tw.TweetFieldReferencedTweets,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("tweet lookup error: %v", err)
	}

	tweetDictionary, ok := tweetResponse.Raw.TweetDictionaries()[tweetID]
	if !ok {
		return nil, fmt.Errorf("tweet lookup error: response does not include tweet with ID %s", tweetID)
	}

	return ParseTweet(tweetDictionary)
}
