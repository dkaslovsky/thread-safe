package thread

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

// Dir constructs the name of the directory where thread files are written
func Dir(topLevelPath string, threadName string) string {
	return filepath.Join(topLevelPath, strings.Replace(threadName, " ", "_", -1))
}

// Load constructs a Thread by loading data from a directory
func Load(path string) (*Thread, error) {
	filePath := filepath.Join(path, fileNameJSON)
	return FromJSON(filePath)
}

// FromJSON constructs a Thread by loading data from a JSON file
func FromJSON(jsonPath string) (*Thread, error) {
	b, err := os.ReadFile(filepath.Clean(jsonPath))
	if err != nil {
		return nil, err
	}

	th := &Thread{}
	return th, json.Unmarshal(b, th)
}

// ToJSON generates and saves a JSON file from a Thread's tweets
func (th *Thread) ToJSON(path string) error {
	err := os.MkdirAll(filepath.Clean(path), 0o750)
	if err != nil {
		return err
	}

	b, bErr := json.Marshal(th)
	if bErr != nil {
		return bErr
	}

	jsonPath := filepath.Clean(filepath.Join(path, fileNameJSON))
	return os.WriteFile(jsonPath, b, 0o600)
}

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

	htmlPath := filepath.Clean(filepath.Join(threadPath, fileNameHTML))

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

// DownloadAttachments saves all media attachments from a Thread's
func (th *Thread) DownloadAttachments(path string) error {
	attachmentPath := filepath.Join(path, dirNameAttachments)
	err := os.MkdirAll(attachmentPath, 0o750)
	if err != nil {
		return err
	}

	for _, tweet := range th.Tweets {
		for _, attachment := range tweet.Attachments {
			attachmentName := attachment.Name(tweet.ID)
			err := attachment.Download(filepath.Clean(filepath.Join(attachmentPath, attachmentName)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadHTMLTemplateFile(threadPath string, templateFile string) (string, error) {
	if templateFile != "" {
		return readFile(templateFile)
	}

	// Try to load default template from file
	defaultTemplateFile := filepath.Clean(filepath.Join(threadPath, "..", fileNameTemplateDefault))
	if _, err := os.Stat(defaultTemplateFile); !os.IsNotExist(err) {
		return readFile(defaultTemplateFile)
	}

	return "", nil
}

func getCSSPath(threadPath string, cssFile string) string {
	if cssFile != "" {
		return filepath.Clean(cssFile)
	}

	// Try to load default CSS file
	defaultCSS := filepath.Clean(filepath.Join(threadPath, "..", fileNameCSSDefault))
	if _, err := os.Stat(defaultCSS); !os.IsNotExist(err) {
		return defaultCSS
	}

	return ""
}

func readFile(path string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
