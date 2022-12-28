package main

import (
	"fmt"
	"os"

	"github.com/dkaslovsky/thread-safe/cmd"
)

// TODO: name vs path everywhere

const (
	name    = "thread-safe"
	version = "0.0.1" // hardcode version for now
)

func main() {
	err := cmd.Run(name, version, os.Args)
	if err != nil {
		fmt.Printf("%s: %v\n", name, err)
		os.Exit(1)
	}
}

// func main() {
// 	token := flag.String("token", "", "twitter API bearer token")                 // TODO: read from env
// 	id := flag.String("id", "", "id of the last tweet in a single-author thread") // TODO: accept id or url
// 	name := flag.String("name", "", "name of thread")                             // TODO: should be arg
// 	pathIn := flag.String("path", "", "path to read/write thread")
// 	write := flag.Bool("write", false, "write data")
// 	attachments := flag.Bool("attachments", false, "download tweet attachments")

// 	devMode := flag.Bool("dev", false, "read pre-saved json file ./test.json") // TODO: remove

// 	flag.Parse()

// 	if *name == "" {
// 		log.Fatal("name is required")
// 	}

// 	// TODO: stat the path to ensure it exists
// 	path := filepath.Join(*pathIn, *name)

// 	// Get thread
// 	var th *Thread
// 	var err error

// 	if *devMode {
// 		th, err = fromFile(filepath.Join(path, "thread.json"))
// 	} else {
// 		th, err = newThreadSaver(*token).thread(*id)
// 	}

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Save thread to filesystem
// 	if *write {
// 		err := os.MkdirAll(path, 0o755)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		werr := th.toFile(path)
// 		if werr != nil {
// 			log.Fatal(werr)
// 		}

// 		if *attachments {
// 			aPath := filepath.Join(path, "attachments")
// 			err := os.MkdirAll(aPath, 0o755)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			serr := th.saveAttachments(aPath)
// 			if serr != nil {
// 				log.Fatal(serr)
// 			}
// 		}
// 	}

// 	// Dump template to console for now
// 	tmpl, err := template.New("thread").Parse(htmlTemplate)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var w io.WriteCloser
// 	if !*write {
// 		w = os.Stdout
// 	} else {
// 		var err error
// 		htmlPath := filepath.Join(path, "thread.html")
// 		w, err = os.Create(htmlPath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	defer func() {
// 		_ = w.Close()
// 	}()

// 	tmplThread, terr := NewTemplateThread(th, *name)
// 	if terr != nil {
// 		log.Fatal(terr)
// 	}
// 	exerr := tmpl.Execute(w, tmplThread)
// 	if exerr != nil {
// 		log.Fatal(exerr)
// 	}
// }
