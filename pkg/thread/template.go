package thread

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TemplateThread represents a top level thread for a template
type TemplateThread struct {
	Name   string          // Name of thread
	Header string          // Thread header information
	Tweets []TemplateTweet // Thread's tweets
}

// TemplateTweet represents a tweet for a template
type TemplateTweet struct {
	Text        string               // Tweet's text contents
	Attachments []TemplateAttachment // Tweet's media attachments
}

// TemplateAttachment represents a tweet's media attachment for a template
type TemplateAttachment struct {
	Path string // Path to the attachment file on the local filesystem
	Ext  string // Attachment's extension
}

// NewTemplateThread constructs a TemplateThread from a thread
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
		Header: th.Metadata(),
		Tweets: tweets,
	}, nil
}

// imageExtensions is a lookup map for identifying image files by extension
var imageExtensions = map[string]struct{}{
	".jpg": {},
	".png": {},
}

// videoExtensions is a lookup map for identifying video files by extension
var videoExtensions = map[string]struct{}{
	".mp4": {},
}

// IsImage evaluates if an attachment is an image file
func (a TemplateAttachment) IsImage() bool {
	_, valid := imageExtensions[a.Ext]
	return valid
}

// IsVideo evaluates if an attachment is a video file
func (a TemplateAttachment) IsVideo() bool {
	_, valid := videoExtensions[a.Ext]
	return valid
}

func loadTemplate(threadDir string, templateFile string, cssFile string) (string, error) {
	html, err := loadHTMLTemplateFile(threadDir, templateFile)
	if err != nil {
		return "", err
	}
	if html == "" {
		html = defaultTemplate
	}

	// Protect against improper formatting if the html template does not provide the "%s" verb
	if !strings.Contains(html, "%s") {
		return html, nil
	}

	return fmt.Sprintf(html, getCSSFile(threadDir, cssFile)), nil
}

const defaultTemplate = `
<head>
<link rel="stylesheet" type="text/css" href="%s" media="screen" />
</head>
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
			<video width="320" height="auto" controls autoplay loop muted><source src=attachments/{{.Path}} type="video/mp4"></video>
			</br></br>
		{{end}}
	{{end}}
{{end}}
`
