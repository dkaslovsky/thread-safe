package twitter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/g8rswimmer/go-twitter/v2"
)

const (
	tweetReferencedTweetTypeRepliedTo = "replied_to"
)

type Tweet struct {
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

func ParseTweet(raw *twitter.TweetDictionary) (*Tweet, error) {
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

	return &Tweet{
		ID:             raw.Tweet.ID,
		ConversationID: raw.Tweet.ConversationID,
		URL:            fmt.Sprintf("https://twitter.com/%s/status/%s", raw.Author.UserName, raw.Tweet.ID),
		Text:           raw.Tweet.Text,
		CreatedAt:      raw.Tweet.CreatedAt,
		AuthorID:       raw.Tweet.AuthorID,
		AuthorName:     raw.Author.Name,
		AuthorHandle:   raw.Author.UserName,
		RepliedToIDs:   repliedToIDs,
		Attachments:    attachments,
	}, nil
}

type Attachment struct {
	MediaKey string `json:"media_key"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

func (a Attachment) Name(tweetID string) string {
	return fmt.Sprintf("tweet=%s-media_key=%s%s", tweetID, a.MediaKey, filepath.Ext(a.URL))
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
	f, fErr := os.Create(path)
	if fErr != nil {
		return fErr
	}
	defer func() {
		_ = f.Close()
	}()

	_, cErr := io.Copy(f, io.LimitReader(resp.Body, 1024*1024))
	if cErr != nil {
		return cErr
	}

	return nil
}