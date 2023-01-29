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

// DirName constructs the name of the directory where thread files are written
func DirName(appDir string, threadName string) string {
	return filepath.Join(appDir, strings.Replace(threadName, " ", "_", -1))
}

func CreateDir(threadDir string) error {
	return os.MkdirAll(filepath.Clean(threadDir), 0o750)
}

// Load constructs a Thread by loading data from a directory
func Load(dir string) (*Thread, error) {
	return FromJSON(filepath.Join(dir, fileNameJSON))
}

// FromJSON constructs a Thread by loading data from a JSON file
func FromJSON(jsonFile string) (*Thread, error) {
	b, err := os.ReadFile(filepath.Clean(jsonFile))
	if err != nil {
		return nil, err
	}

	th := &Thread{}
	return th, json.Unmarshal(b, th)
}

// ToJSON generates and saves a JSON file from a Thread's tweets
func (th *Thread) ToJSON(threadDir string) error {
	b, err := json.Marshal(th)
	if err != nil {
		return err
	}

	jsonPath := filepath.Clean(filepath.Join(threadDir, fileNameJSON))
	return os.WriteFile(jsonPath, b, 0o600)
}

// ToHTML generates and saves an HTML file from a thread using default or provided template and CSS files
func (th *Thread) ToHTML(threadDir string, templateFile string, cssFile string) error {
	htmlTemplate, err := loadTemplate(threadDir, templateFile, cssFile)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	tmpl, tErr := template.New("thread").Parse(htmlTemplate)
	if tErr != nil {
		return fmt.Errorf("failed to parse template: %w", tErr)
	}

	htmlFile := filepath.Clean(filepath.Join(threadDir, fileNameHTML))

	tmplThread, ttErr := NewTemplateThread(th)
	if ttErr != nil {
		log.Fatal(ttErr)
	}
	f, fErr := os.Create(htmlFile)
	if fErr != nil {
		return fErr
	}
	defer func() {
		_ = f.Close()
	}()

	return tmpl.Execute(f, tmplThread)
}

// DownloadAttachments saves all media attachments from a Thread's
func (th *Thread) DownloadAttachments(threadDir string) error {
	attachmentDir := filepath.Join(threadDir, dirNameAttachments)
	err := os.MkdirAll(attachmentDir, 0o750)
	if err != nil {
		return err
	}

	for _, tweet := range th.Tweets {
		for _, attachment := range tweet.Attachments {
			err := attachment.Download(filepath.Clean(filepath.Join(attachmentDir, attachment.Name(tweet.ID))))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadHTMLTemplateFile(threadDir string, templateFile string) (string, error) {
	if templateFile != "" {
		return readFile(templateFile)
	}

	// Try to load default template from file
	defaultTemplateFile := filepath.Clean(filepath.Join(threadDir, "..", fileNameTemplateDefault))
	if _, err := os.Stat(defaultTemplateFile); !os.IsNotExist(err) {
		return readFile(defaultTemplateFile)
	}

	return "", nil
}

func getCSSFile(threadDir string, cssFile string) string {
	if cssFile != "" {
		return filepath.Clean(cssFile)
	}

	// Try to load default CSS file
	defaultCSSFile := filepath.Clean(filepath.Join(threadDir, "..", fileNameCSSDefault))
	if _, err := os.Stat(defaultCSSFile); !os.IsNotExist(err) {
		return defaultCSSFile
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
