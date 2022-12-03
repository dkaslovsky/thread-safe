package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/g8rswimmer/go-twitter/v2"
)

type Client interface {
	tweetLookup(id string, opts twitter.TweetLookupOpts) (*tweet, error)
}

func newClient(token string) Client {
	c := &twitter.Client{
		Authorizer: authorize{
			Token: token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	return &client{c}
}

type client struct {
	*twitter.Client
}

func (c *client) tweetLookup(tweetID string, opts twitter.TweetLookupOpts) (*tweet, error) {
	tweetResponse, err := c.TweetLookup(context.Background(), []string{tweetID}, opts)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup error: %v", err)
	}

	tweetDictionary, ok := tweetResponse.Raw.TweetDictionaries()[tweetID]
	if !ok {
		return nil, fmt.Errorf("tweet lookup error: response does not include tweet with ID %s", tweetID)
	}

	return parseTweet(tweetDictionary)
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}
