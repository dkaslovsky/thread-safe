package thread

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	htmlFileName = "thread.html"
)

func (th *Thread) ToHTML(path string, css string) error {
	tmpl, err := template.New("thread").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	htmlPath := filepath.Clean(filepath.Join(path, htmlFileName))
	cssPath := ""
	if css != "" {
		cssPath = filepath.Clean(css)
	}

	tmplThread, tErr := NewTemplateThread(th, cssPath)
	if tErr != nil {
		log.Fatal(tErr)
	}
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

func NewTemplateThread(th *Thread, cssPath string) (TemplateThread, error) {
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
		CSS:    cssPath,
	}, nil
}

type TemplateThread struct {
	Name   string
	Header string
	Tweets []TemplateTweet
	CSS    string
}

func (t TemplateThread) HasCSS() bool {
	return t.CSS != ""
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

const htmlTemplate = `
{{if .HasCSS}}
	<head>
	<link rel="stylesheet" type="text/css" href="../thread-safe.css" media="screen" />
	</head>
{{end}}
<h1>{{.Name}}</h1>
<div class="text"><pre>{{.Header}}</pre></div>
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
`
