# thread-safe
Keep your favorite Twitter threads safe by downloading a local copy


## Overview
`thread-safe` is a simple CLI for saving a local copy of a Twitter thread.
Specifically, thread-safe writes a .json file containing the thread data, downloads all of the thread's media attachments (images, videos), and generates an .html file for displaying the thread in a browser.

`thread-safe` is designed to
* save a local copy of information-rich Twitter threads
* eliminate the need to use an external or third-party app to which you need to grant access to your Twitter account
* eliminate the need to reply to a tweet to "unroll" or otherwise save copy of the contents
* ignore replies from users other than the thread's author

Therefore, thread-safe's definition of a thread is intentionally limited to: _a consecutive series of tweets authored by the same Twitter user_. 

## Example
TODO: Find a good, non-controversial thread with media attachments as an example...

## Installation
TODO

## Usage
`thread-safe` is lightweight and simple to use. To save a thread, two items are needed:
* A valid Twitter API bearer token (see Twitter's [documentation](https://developer.twitter.com/en/docs/authentication/oauth-2-0/bearer-tokens))
* The URL or ID of the **last** tweet in the thread

A tweet's URL is typically of the form `https://twitter.com/<username>/status/<ID>?<params>`.
The entire URL or simply the numeric `<ID>` portion of the URL can be provided to identify a tweet.

While it is a bit annoying to have to identify the last tweet in the thread rather than the more natural first tweet, limitations of the Twitter API make this unavoidable for `thread-safe`'s workflow.

</br>

### Top Level Usage
```
$ thread-safe --help
thread-safe saves a local copy of a Twitter thread

Usage:
  thread-safe [flags]
  thread-safe [command]

Available Commands:
  thread  saves thread content and generates a local html file
  html    regenerates an html file from a previously saved thread

Flags:
  -h, --help	 help for thread-safe
  -v, --version	 version for thread-safe

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (current directory if unset)
  THREAD_SAFE_TOKEN	bearer token for Twitter API

Use "thread-safe [command] --help" for more information about a command
```
Note that the Twitter API bearer token must be provided via the `THREAD_SAFE_TOKEN` environment variable and that files will be written to either the directory specified by `THREAD_SAFE_PATH` or the current directory if the environment variable is not set.

</br>

### Subcommands
`thread`
```
thread saves thread content and generates a local html file

Usage:
  thread-safe thread [flags] <name> <last-tweet>

Args:
  name           string  name to use for the thread
  last-tweet     string  URL or ID of the last tweet in a single-author thread

Flags:
  -c, --css             string  optional path to CSS file
  -t, --template        string  optional path to template file
      --no-attachments          do not download attachments

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (current directory if unset)
  THREAD_SAFE_TOKEN	bearer token for Twitter API
```

`html`
```
html regenerates an html file from a previously saved thread

Usage:
  thread-safe html [flags] <name>

Args:
  name  string  name given to the thread

Flags:
  -c, --css       string  optional path to CSS file
  -t, --template  string  optional path to template file

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (current directory if unset)
  THREAD_SAFE_TOKEN	bearer token for Twitter API
```

#### **Custom CSS**
The `thread` and `html` subcommands support providing an optional path to a CSS file to be linked as an external stylesheet in the generated HTML.

#### **Custom Templates**
The `thread` and `html` subcommands also support providing an optional path to a file containing an HTML template to be used in place of `thread-safe`'s default template. The contents of a provided template file must be parsable by Golang's [(*Template).Parse](https://pkg.go.dev/text/template#Template.Parse) function.

The template must make use of the following objects:

The top level `TemplateThread` is defined by
```go
type TemplateThread struct {
	Name   string          // Name of the thread
	Header string          // Thread header information
	CSS    string          // Path to a custom CSS file
	Tweets []TemplateTweet // Thread's tweets
}

func (TemplateThread) HasCSS() bool
```
The nested `TemplateTweet` object is defined by
```go
type TemplateTweet struct {
	Text        string               // Tweet's text contents
	Attachments []TemplateAttachment // Tweet's media attachments
}
```
with the `TemplateAttachment` object defined by
```go
type TemplateAttachment struct {
	Path string // Path to the attachment file on the local filesystem
	Ext  string // Attachment's extension (.jpg, .mpe4
}

func (TemplateAttachment) IsImage() bool

func (a TemplateAttachment) IsVideo() bool
```
