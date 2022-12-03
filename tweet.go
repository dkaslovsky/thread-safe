package main

import "github.com/g8rswimmer/go-twitter/v2"

// TODO: unexported fields for now
type tweet struct {
	ID             string   `json:"id"`
	ConversationID string   `json:"conversation_id"`
	URL            string   `json:"url"`
	Text           string   `json:"text"`
	CreatedAt      string   `json:"created_at"`
	AuthorID       string   `json:"author_id"`
	AuthorName     string   `json:"author_name"`
	AuthorHandle   string   `json:"author_handle"`
	RepliedToIDs   []string `json:"replied_to_ids"`
	// Attachments  []Attachment // TODO
}

func parseTweet(raw *twitter.TweetDictionary) (*tweet, error) {
	repliedToIDs := []string{}
	for _, ref := range raw.Tweet.ReferencedTweets {
		if ref.Type == "replied_to" {
			repliedToIDs = append(repliedToIDs, ref.ID)
		}
	}

	return &tweet{
		ID:             raw.Tweet.ID,
		ConversationID: raw.Tweet.ConversationID,
		URL:            "", // TODO
		Text:           raw.Tweet.Text,
		CreatedAt:      raw.Tweet.CreatedAt,
		AuthorID:       raw.Tweet.AuthorID,
		AuthorName:     raw.Author.Name,
		AuthorHandle:   raw.Author.UserName,
		RepliedToIDs:   repliedToIDs,
	}, nil
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
