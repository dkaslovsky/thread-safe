package thread

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const htmlFileName = "thread.html"

// TODO:
// - improve template formatting
// - support CSS
// - possibly support user-provided CSS and/or template
const htmlTemplate = `
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

func (th *Thread) ToHTML(path string) error {
	tmpl, err := template.New("thread").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	tmplThread, tErr := NewTemplateThread(th)
	if tErr != nil {
		log.Fatal(tErr)
	}

	htmlPath := filepath.Clean(filepath.Join(path, htmlFileName))
	f, fErr := os.Create(htmlPath)
	if fErr != nil {
		return fErr
	}
	defer func() {
		_ = f.Close()
	}()

	return tmpl.Execute(f, tmplThread)
}

func (th *Thread) Header() string {
	if th.Len() == 0 {
		return ""
	}
	first := th.Tweets[0]
	headerStrs := []string{
		fmt.Sprintf("URL: \t\t\t%s", first.URL),
		fmt.Sprintf("Author Name: \t\t%s", first.AuthorName),
		fmt.Sprintf("Author Handle: \t\t%s", first.AuthorHandle),
		fmt.Sprintf("Conversation ID: \t%s", first.ConversationID),
	}
	return strings.Join(headerStrs, "\n")
}

func NewTemplateThread(th *Thread) (TemplateThread, error) {
	threadLen := th.Len()
	tweets := []TemplateTweet{}
	for i, tweet := range th.Tweets {
		attachments := []TemplateAttachment{}
		for _, attachment := range tweet.Attachments {
			path := attachment.Name(tweet.ID)
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
		Name:   th.Name,
		Header: th.Header(),
		Tweets: tweets,
	}, nil
}

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
