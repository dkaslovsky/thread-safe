package main

import (
	"fmt"
	"path/filepath"
)

var htmlTemplate = `
<h1>{{.Name}}</h1>
<div class="text"><pre>{{.Header}}</pre></div>
<ul>
    {{range .Tweets}}
		<h3>{{.Text}}</h3>
		</br></br>
		{{range .Attachments}}
			{{.HTML}}
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
	HTML string
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

			var html string
			switch filepath.Ext(name) {
			case ".jpg":
				html = `<img width="320" height="auto" src=attachments/%s>`
			case ".mp4":
				html = `<video width="320" height="240" controls><source src=attachments/%s type="video/mp4"></video>`
			}

			if html == "" {
				return TemplateThread{}, fmt.Errorf("unknown attachment type %s", filepath.Ext(name))
			}

			attachments = append(attachments, TemplateAttachment{
				HTML: fmt.Sprintf(html+"</br></br>", name),
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
