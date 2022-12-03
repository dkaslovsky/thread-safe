package main

import (
	"fmt"
	"strings"
)

type thread struct {
	tweets []*tweet
}

func (t *thread) Len() int {
	return len(t.tweets)
}

func (t *thread) String() string {
	tweetStrs := []string{}
	threadLen := t.Len()
	for i, tweet := range t.tweets {
		tweetStr := fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text)
		tweetStrs = append(tweetStrs, tweetStr)
	}
	return strings.Join(tweetStrs, "\n---\n")
}
