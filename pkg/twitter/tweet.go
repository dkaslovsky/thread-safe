package twitter

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	tw "github.com/g8rswimmer/go-twitter/v2"
)

const (
	// tweetReferencedTweetTypeRepliedTo is the field to use for following a thread's response chain
	tweetReferencedTweetTypeRepliedTo = "replied_to"
)

// Tweet represents a Twitter tweet
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

// ParseTweet constructs a Tweet from the data returned by querying the Twitter API
func ParseTweet(raw *tw.TweetDictionary) (*Tweet, error) {
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
			continue
		}

		// Add the attachment variant with the largest bit rate
		type variantAttachment struct {
			url     string
			bitRate int
		}
		curVariant := variantAttachment{}
		for _, variant := range attachement.Variants {
			if variant.URL == "" {
				continue
			}
			if variant.BitRate > curVariant.bitRate {
				curVariant.url = variant.URL
				curVariant.bitRate = variant.BitRate
			}
		}
		if curVariant.url != "" {
			attachments = append(attachments, Attachment{
				MediaKey: attachement.Key,
				Type:     attachement.Type,
				URL:      curVariant.url,
			})
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

// Attachment represents a media file attached to a Tweet
type Attachment struct {
	MediaKey string `json:"media_key"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

// Name constructs the file name to use for saving an Attachment
func (a Attachment) Name(tweetID string) string {
	// Clean the file extension by removing any invalid params
	ext := strings.SplitN(filepath.Ext(a.URL), "?", 2)[0]
	return fmt.Sprintf("tweet=%s-media_key=%s%s", tweetID, a.MediaKey, ext)
}

// Download saves an Attachment as a file
func (a Attachment) Download(path string) error {
	if u, err := url.ParseRequestURI(a.URL); !(err == nil && u.Scheme != "" && u.Host != "") {
		return fmt.Errorf("invalid attachment URL %s for media_key %s", a.URL, a.MediaKey)
	}
	resp, err := http.Get(a.URL)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download of %s failed with status code: %d", a.URL, resp.StatusCode)
	}

	f, fErr := os.Create(filepath.Clean(path))
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
