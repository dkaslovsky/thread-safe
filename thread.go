package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type thread struct {
	items []*tweet
}

func (t *thread) tweets() []*tweet {
	return t.items
}

func (t *thread) len() int {
	return len(t.items)
}

func (t *thread) String() string {
	tweetStrs := []string{t.header()}
	threadLen := t.len()
	for i, tweet := range t.tweets() {
		tweetStr := fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text)
		tweetStrs = append(tweetStrs, tweetStr)
	}
	return strings.Join(tweetStrs, "\n---\n")
}

func (t *thread) header() string {
	if t.len() == 0 {
		return ""
	}
	first := t.tweets()[0]
	headerStrs := []string{
		fmt.Sprintf("URL: \t\t\t%s", first.URL), // TODO: sanitize HTML here or in template
		fmt.Sprintf("Author Name: \t\t%s", first.AuthorName),
		fmt.Sprintf("Author Handle: \t\t%s", first.AuthorHandle),
		fmt.Sprintf("Conversation ID: \t%s", first.ConversationID),
	}
	return strings.Join(headerStrs, "\n")
}

func (t *thread) toFile(path string) error {
	b, err := json.Marshal(t.items)
	if err != nil {
		return err
	}
	filePath := filepath.Join(path, "thread.json")
	fmt.Printf("file path = %s\n", filePath)
	return os.WriteFile(filePath, b, 0o755)
}

func fromFile(path string) (*thread, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := &thread{}
	return t, json.Unmarshal(b, &(t.items))
}
