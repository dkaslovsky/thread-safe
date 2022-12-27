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
		{{range .Attachments}}
			<li><img src=attachments/{{.Path}}></li>
    	{{end}}
    {{end}}
</ul>
`

type TemplateThread struct {
	Name   string
	Header string
	Tweets []TemplateTweet
}

type TemplateTweet struct {
	Text        string
	Attachments []TemplateAttachment
}

type TemplateAttachment struct {
	Path string
}

func NewTemplateThread(t *thread, name string) (TemplateThread, error) {
	threadLen := t.len()
	tweets := []TemplateTweet{}
	for i, tweet := range t.tweets() {
		attachments := []TemplateAttachment{}
		for i := 0; i < len(tweet.Attachments); i++ {
			name, nerr := tweet.AttachmentName(i)
			if nerr != nil {
				return TemplateThread{}, nerr
			}
			attachments = append(attachments, TemplateAttachment{
				Path: name,
			})
		}
		tt := TemplateTweet{
			Text:        fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text),
			Attachments: attachments,
		}
		tweets = append(tweets, tt)
	}

	return TemplateThread{
		Name:   name,
		Header: t.header(),
		Tweets: tweets,
	}, nil
}
