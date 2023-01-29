package thread

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const (
	// dirNameAttachments is the name used for the directory where attachment files are saved
	dirNameAttachments = "attachments"

	// fileNameCSSDefault is the CSS file used if it exists and no other CSS file is specified
	fileNameCSSDefault = "thread-safe.css"
	// fileNameTemplateDefault is the template file used if it exists and no other template file is specified
	fileNameTemplateDefault = "thread-safe.tmpl"
	// fileNameHTML is the name used for the generated HTML file
	fileNameHTML = "thread.html"
	// fileNameJSON is the name used for the generated JSON file
	fileNameJSON = "thread.json"
)

// FromJSON constructs a Thread by loading data from a JSON file
func FromJSON(appDir string, threadName string) (*Thread, error) {
	dir := NewDirectory(appDir, threadName)

	if !dir.Exists() {
		return nil, fmt.Errorf("%s not found", dir)
	}

	b, err := os.ReadFile(dir.Join(fileNameJSON))
	if err != nil {
		return nil, err
	}

	th := Thread{}
	jErr := json.Unmarshal(b, &th)
	if jErr != nil {
		return nil, jErr
	}

	th.Dir = dir
	return &th, nil
}

// ToJSON generates and saves a JSON file from a Thread's tweets
func (th *Thread) ToJSON() error {
	b, err := json.Marshal(th)
	if err != nil {
		return err
	}

	return os.WriteFile(th.Dir.Join(fileNameJSON), b, 0o600)
}

// ToHTML generates and saves an HTML file from a thread using default or provided template and CSS files
func (th *Thread) ToHTML(templateFile string, cssFile string) error {
	htmlTemplate, err := loadTemplate(th.Dir, templateFile, cssFile)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	tmpl, tErr := template.New("thread").Parse(htmlTemplate)
	if tErr != nil {
		return fmt.Errorf("failed to parse template: %w", tErr)
	}

	f, fErr := os.Create(th.Dir.Join(fileNameHTML))
	if fErr != nil {
		return fmt.Errorf("failed to open HTML file: %w", fErr)
	}
	defer func() {
		_ = f.Close()
	}()

	eErr := tmpl.Execute(f, NewTemplateThread(th))
	if eErr != nil {
		return fmt.Errorf("failed to execute template: %w", eErr)
	}
	return nil
}

// DownloadAttachments saves all media attachments from a Thread's
func (th *Thread) DownloadAttachments() error {
	attachmentDir := NewDirectory(th.Dir.Join(dirNameAttachments), "")
	err := attachmentDir.Create()
	if err != nil {
		return err
	}

	for _, tweet := range th.Tweets {
		for _, attachment := range tweet.Attachments {
			err := attachment.Download(attachmentDir.Join(attachment.Name(tweet.ID)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadHTMLTemplateFile(threadDir *Directory, templateFile string) (string, error) {
	if templateFile != "" {
		return readFile(templateFile)
	}

	// Try to load default template from file
	if defaultFile, exists := threadDir.SubDir("..", fileNameTemplateDefault); exists {
		return readFile(defaultFile)
	}

	return "", nil
}

func getCSSFile(threadDir *Directory, cssFile string) string {
	if cssFile != "" {
		return filepath.Clean(cssFile)
	}

	// Try to load default CSS file
	if defaultFile, exists := threadDir.SubDir("..", fileNameCSSDefault); exists {
		return defaultFile
	}

	return ""
}

func readFile(fileName string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
