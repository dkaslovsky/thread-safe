# thread-safe
Keep your favorite Twitter threads safe by downloading a local copy

## Overview
`thread-safe` is a simple CLI for saving a local copy of a Twitter thread.

Specifically, `thread-safe` generates an HTML file containing all of a thread's contents including each tweet's text, links, and media attachments (images, videos). This file, all attachments, and a JSON data file are saved to the local filesystem and the HTML can be used to display the thread locally in a browser at any time.

By using a dedicated directory for all generated files, `thread-safe` can be used to maintain a local library of saved threads. Thread names are specified by the user as CLI arguments and standard commandline tooling (e.g., `grep`, `find`, `fzf`, etc) can be used to search the library for saved content.

`thread-safe` is designed to
* Save a local copy of informative Twitter threads
* Eliminate the need to use an external or third-party app to which you need to grant access to your Twitter account
* Eliminate the need to reply to a tweet to "unroll" or otherwise save the thread's contents
* Ignore replies from users other than the thread's author

Therefore, `thread-safe`'s definition of a thread is intentionally limited to
> _a consecutive series of tweets authored by the same Twitter user_. 

## Example
To demonstrate typical usage, we've identified a Twitter thread of [Nathan MacKinnon hockey highlights](https://twitter.com/Avalanche/status/969990878944149504) from 2018 that we simply must preserve with a local copy. This thread is great not only for its content but also because it contains both images and video. The thread contains eight tweets before any non-author replies and the URL of the last tweet URL is `https://twitter.com/Avalanche/status/969990907490484225`.

We'll name our local copy of this thread `Nathan MacKinnon 2018` and save it by running
```
$ thread-safe save "Nathan MacKinnon 2018" "https://twitter.com/Avalanche/status/969990907490484225"
```
or, equivalently,
```
$ thread-safe save "Nathan MacKinnon 2018" 969990907490484225
```

We now have a `thread.html` file in the `$THREAD_SAFE_PATH/Nathan_MacKinnon_2018` directory with the all of the thread's contents:

https://user-images.githubusercontent.com/20505301/210184312-dcc6be97-3d1f-4fff-ba0e-36cd00f52add.mov

Note that we ran `thread-safe` with the default HTML template and no CSS. If we later wish to regenerate `thread.html` with a specified template and CSS, we can run
```
$ thread-safe regen --template path/to/template --css path/to/css "Nathan MacKinnon 2018"
```
and the file will be rewritten using the target template and CSS files. We also could have provided the `--template` and `--css` flags to the original `thread-safe save` command.

All resulting thread files can be found in this repository's [examples](examples) directory.

## Installation
TODO

## Usage
`thread-safe` is lightweight and simple to use. To save a thread, two items are needed:
* A valid [Twitter API bearer token](https://developer.twitter.com/en/docs/authentication/oauth-2-0/bearer-tokens)
* The URL or ID of the **last** tweet in the thread

A tweet's URL is typically of the form
>`https://twitter.com/<username>/status/<ID>?<parameters>`

The entire URL or simply the numeric `<ID>` portion of the URL can be provided as an argument to specify a tweet.

While it is inconvenient to have to identify the last tweet in the thread rather than the more natural first tweet, limitations of the Twitter API make this unavoidable for `thread-safe`'s workflow.

</br>

### Top Level Usage
```
$ thread-safe --help
'thread-safe' saves a local copy of a Twitter thread

Usage:
  thread-safe [flags]
  thread-safe [command]

Available Commands:
  save     saves thread content and generates a local html file
  regen    regenerates an html file from a previously saved thread

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
* `save`: save thread data and generate HTML for local browsing
```
$ thread-safe save --help
'save' saves thread content and generates a local html file

Usage:
  thread-safe save [flags] <name> <last-tweet>

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

* `regen`: reprocess saved thread data using an updated template or CSS
```
$ thread-safe regen --help
'regen' regenerates an html file from a previously saved thread

Usage:
  thread-safe regen [flags] <name>

Args:
  name  string  name given to the thread

Flags:
  -c, --css       string  optional path to CSS file
  -t, --template  string  optional path to template file

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (current directory if unset)
  THREAD_SAFE_TOKEN	bearer token for Twitter API
```
</br>

### Custom CSS
The `save` and `regen` subcommands support providing an optional path to a CSS file to be linked as an external stylesheet in the generated HTML.

</br>

### Custom Templates
The `save` and `regen` subcommands also support providing an optional path to a file containing an HTML template to be used in place of `thread-safe`'s default template. The contents of a provided template file must be parsable by the Go [(*Template).Parse()](https://pkg.go.dev/text/template#Template.Parse) function.

The template must make use of the following objects:

* The top level `TemplateThread` object defined by
```go
type TemplateThread struct {
	Name   string          // Name of thread
	Header string          // Thread header information
	CSS    string          // Path to custom CSS file
	Tweets []TemplateTweet // Thread's tweets
}

func (TemplateThread) HasCSS() bool
```
* The nested `TemplateTweet` object defined by
```go
type TemplateTweet struct {
	Text        string               // Tweet's text content
	Attachments []TemplateAttachment // Tweet's media attachments
}
```
* The `TemplateAttachment` object defined by
```go
type TemplateAttachment struct {
	Path string // Path to the attachment file on the local filesystem
	Ext  string // Attachment's extension (.jpg, .mp4)
}

func (TemplateAttachment) IsImage() bool

func (TemplateAttachment) IsVideo() bool
```

## License
thread-safe is released under the [MIT License](./LICENSE).
Dependency licenses are available in this repository's [CREDITS](./CREDITS) file.
