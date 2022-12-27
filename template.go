package main

import (
	"fmt"
)

var htmlTemplate = `
<h1>{{.Name}}</h1>
<div class="text"><pre>{{.Header}}</pre></div>
<ul>
    {{range .Tweets}}
		<li>{{.Text}}</li>
    {{end}}
</ul>
`

type TemplateThread struct {
	Name   string
	Header string
	Tweets []TemplateTweet
}

type TemplateTweet struct {
	Text string
}

func NewTemplateThread(t *thread, name string) TemplateThread {
	threadLen := t.len()
	tweets := []TemplateTweet{}
	for i, tweet := range t.tweets() {
		tweets = append(tweets, TemplateTweet{
			Text: fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text),
		})
	}

	return TemplateThread{
		Name:   name,
		Header: t.header(),
		Tweets: tweets,
	}
}
