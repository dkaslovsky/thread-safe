package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// maxThreadLen is the maximum number of tweets allowed to be fetched for constructing a thread
const maxThreadLen = 100

type threadSaver struct {
	client Client
}

func newThreadSaver(token string) *threadSaver {
	return &threadSaver{
		client: newClient(token),
	}
}

func (ts *threadSaver) thread(lastID string) (*thread, error) {
	tweets, err := ts.walkTweets(lastID, maxThreadLen)
	if err != nil {
		return nil, err
	}

	// Tweets are fetched from last to first so reverse the order
	reverseSlice(tweets)

	return &thread{tweets}, nil
}

func (ts *threadSaver) walkTweets(id string, limit int) ([]*tweet, error) {
	tweets := []*tweet{}

	nextID := id
	conversationID := ""
	authorID := ""

	for i := 0; i < limit; i++ {
		tweet, err := ts.client.tweetLookup(nextID)
		if err != nil {
			return nil, err
		}

		// Save the conversationID and authorID
		if i == 0 {
			conversationID = tweet.ConversationID
			authorID = tweet.AuthorID
		}

		// A change in conversationID or authorID indicates the end of the current thread
		if tweet.ConversationID != conversationID || tweet.AuthorID != authorID {
			return tweets, nil
		}

		tweets = append(tweets, tweet)

		switch len(tweet.RepliedToIDs) {
		case 1: // Next ID for lookup
			nextID = tweet.RepliedToIDs[0]
		case 0: // Top of thread has been reached
			return tweets, nil
		default: // Error on multiple replied_to IDs
			return nil, fmt.Errorf("cannot follow tweet %s with multiple replied_to IDs", tweet.ID)
		}
	}

	// Limit reached
	return nil, fmt.Errorf("exceeded maximum number of tweets to fetch [%d]", limit)
}

func reverseSlice[T any](s []T) {
	first, last := 0, len(s)-1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}

func main() {
	token := flag.String("token", "", "twitter API bearer token")                 // TODO: read from env
	id := flag.String("id", "", "id of the last tweet in a single-author thread") // TODO: accept id or url
	name := flag.String("name", "", "name of thread")                             // TODO: should be arg
	writePath := flag.String("path", "", "path to write thread to file")
	attachments := flag.Bool("attachments", false, "download tweet attachments")

	devMode := flag.Bool("dev", false, "read pre-saved json file ./test.json") // TODO: remove
	flag.Parse()

	if *name == "" {
		log.Fatal("name is required")
	}

	// TODO: stat the path to ensure it exists
	path := filepath.Join(*writePath, *name)

	// Get thread
	var th *thread
	var err error

	if *devMode {
		th, err = fromFile("./test.json")
	} else {
		th, err = newThreadSaver(*token).thread(*id)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Save thread to filesystem
	if *writePath != "" {
		err := os.MkdirAll(path, 0o755)
		if err != nil {
			log.Fatal(err)
		}
		werr := th.toFile(path)
		if werr != nil {
			log.Fatal(werr)
		}

		if *attachments {
			aPath := filepath.Join(path, "attachments")
			err := os.MkdirAll(aPath, 0o755)
			if err != nil {
				log.Fatal(err)
			}
			serr := th.saveAttachments(aPath)
			if serr != nil {
				log.Fatal(serr)
			}
		}
	}

	// Dump template to console for now
	tmpl, err := template.New("thread").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var w io.WriteCloser
	if *writePath == "" {
		w = os.Stdout
	} else {
		var err error
		htmlPath := filepath.Join(path, "thread.html")
		w, err = os.Create(htmlPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func() {
		_ = w.Close()
	}()

	terr := tmpl.Execute(w, NewTemplateThread(th, *name))
	if terr != nil {
		log.Fatal(terr)
	}
}
