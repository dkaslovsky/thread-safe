package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/g8rswimmer/go-twitter/v2"
)

// TODO: unexported fields for now

type tweet struct {
	ID             string       `json:"id"`
	ConversationID string       `json:"conversation_id"`
	URL            string       `json:"url"`
	Text           string       `json:"text"`
	CreatedAt      string       `json:"created_at"`
	AuthorID       string       `json:"author_id"`
	AuthorName     string       `json:"author_name"`
	AuthorHandle   string       `json:"author_handle"`
	RepliedToIDs   []string     `json:"replied_to_ids"`
	Attachments    []Attachment `json:"attachments"`
}

const tweetReferencedTweetTypeRepliedTo = "replied_to"

func parseTweet(raw *twitter.TweetDictionary) (*tweet, error) {
	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", raw.Author.UserName, raw.Tweet.ID)

	repliedToIDs := []string{}
	for _, ref := range raw.Tweet.ReferencedTweets {
		if ref.Type == tweetReferencedTweetTypeRepliedTo {
			repliedToIDs = append(repliedToIDs, ref.ID)
		}
	}

	attachments := []Attachment{}
	for _, attachement := range raw.AttachmentMedia {
		if attachement.URL != "" {
			attachments = append(attachments, Attachment{
				MediaKey: attachement.Key,
				Type:     attachement.Type,
				URL:      attachement.URL,
			})
		}
		for _, variant := range attachement.Variants {
			if variant.URL != "" {
				attachments = append(attachments, Attachment{
					MediaKey: attachement.Key,
					Type:     attachement.Type,
					URL:      variant.URL,
				})
			}
		}
	}

	return &tweet{
		ID:             raw.Tweet.ID,
		ConversationID: raw.Tweet.ConversationID,
		URL:            tweetURL,
		Text:           raw.Tweet.Text,
		CreatedAt:      raw.Tweet.CreatedAt,
		AuthorID:       raw.Tweet.AuthorID,
		AuthorName:     raw.Author.Name,
		AuthorHandle:   raw.Author.UserName,
		RepliedToIDs:   repliedToIDs,
		Attachments:    attachments,
	}, nil
}

func (t *tweet) AttachmentName(idx int) (string, error) {
	if idx >= len(t.Attachments) {
		return "", fmt.Errorf("invalid attachment index %d for tweet with %d attachments", idx, len(t.Attachments))
	}
	a := t.Attachments[idx]
	return fmt.Sprintf("tweet=%s-media_key=%s%s", t.ID, a.MediaKey, filepath.Ext(a.URL)), nil
}

type Attachment struct {
	MediaKey string `json:"media_key"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

func (a Attachment) Download(path string) error {
	// TODO: sanitize or check url
	resp, err := http.Get(a.URL)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Download of [%s] failed with status code: %d", a.URL, resp.StatusCode)
	}

	// TODO: sanitize or check path
	f, ferr := os.Create(path)
	if ferr != nil {
		return ferr
	}
	defer func() {
		_ = f.Close()
	}()

	_, cerr := io.Copy(f, io.LimitReader(resp.Body, 1024*1024))
	if cerr != nil {
		return err
	}

	return nil
}
