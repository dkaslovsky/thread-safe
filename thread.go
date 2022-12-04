package main

import (
	"fmt"
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
	tweetStrs := []string{}
	threadLen := t.len()
	for i, tweet := range t.tweets() {
		tweetStr := fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text)
		tweetStrs = append(tweetStrs, tweetStr)
	}
	return strings.Join(tweetStrs, "\n---\n")
}
