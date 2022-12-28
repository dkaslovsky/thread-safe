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
			{{if .IsImage}}
				<img width="320" height="auto" src=attachments/{{.Path}}>
				</br></br>
			{{end}}
			{{if .IsVideo}}
				<video width="320" height="240" controls autoplay loop muted><source src=attachments/{{.Path}} type="video/mp4"></video>
				</br></br>
			{{end}}
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
	Ext  string
}

func (a TemplateAttachment) IsImage() bool {
	return a.Ext == ".jpg"
}

func (a TemplateAttachment) IsVideo() bool {
	return a.Ext == ".mp4"
}

func NewTemplateThread(t *thread, name string) (TemplateThread, error) {
	threadLen := t.len()
	tweets := []TemplateTweet{}
	for i, tweet := range t.tweets() {
		attachments := []TemplateAttachment{}
		for i := 0; i < len(tweet.Attachments); i++ {
			path, err := tweet.AttachmentName(i)
			if err != nil {
				return TemplateThread{}, err
			}
			attachments = append(attachments, TemplateAttachment{
				Path: path,
				Ext:  filepath.Ext(path),
			})
		}

		tweets = append(tweets, TemplateTweet{
			Text:        fmt.Sprintf("[%d/%d] %s", i+1, threadLen, tweet.Text),
			Attachments: attachments,
		})
	}

	return TemplateThread{
		Name:   name,
		Header: t.header(),
		Tweets: tweets,
	}, nil
}
