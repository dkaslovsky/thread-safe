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
	// htmlFileName is the name used for the generated HTML file
	htmlFileName = "thread.html"

	// defaultCSSFileName is the CSS file used if it exists and no other CSS file is specified
	defaultCSSFileName = "thread-safe.css"
)

// ToHTML generates and saves an HTML file from a thread using default or provided template and CSS files
func (th *Thread) ToHTML(threadPath string, templateFile string, cssFile string) error {
	htmlTemplate, err := loadTemplate(threadPath, templateFile, cssFile)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	tmpl, tErr := template.New("thread").Parse(htmlTemplate)
	if tErr != nil {
		return fmt.Errorf("failed to parse template: %w", tErr)
	}

	htmlPath := filepath.Clean(filepath.Join(threadPath, htmlFileName))

	tmplThread, ttErr := NewTemplateThread(th)
	if ttErr != nil {
		log.Fatal(ttErr)
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

// Header returns a string with thread metadata
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
		Header: th.Header(),
		Tweets: tweets,
	}, nil
}

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

func loadTemplate(threadPath string, templateFile string, cssFile string) (string, error) {
	html, err := loadHTMLTemplateFile(templateFile)
	if err != nil {
		return "", err
	}

	// Ignore CSS if the html template does not provide the "%s" format verb for it to be included
	if !strings.Contains(html, "%s") {
		return html, nil
	}

	return fmt.Sprintf(html, getCSSPath(threadPath, cssFile)), nil
}

func loadHTMLTemplateFile(path string) (string, error) {
	if path == "" {
		return defaultTemplate, nil
	}
	templatePath := filepath.Clean(path)
	b, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getCSSPath(threadPath string, cssFile string) string {
	if cssFile != "" {
		return filepath.Clean(cssFile)
	}

	// Try to load global CSS file
	defaultCSS := filepath.Clean(filepath.Join(threadPath, "..", defaultCSSFileName))
	if _, err := os.Stat(defaultCSS); !os.IsNotExist(err) {
		return defaultCSS
	}

	return ""
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
			<video width="320" height="240" controls autoplay loop muted><source src=attachments/{{.Path}} type="video/mp4"></video>
			</br></br>
		{{end}}
	{{end}}
{{end}}
`
