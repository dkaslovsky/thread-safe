# thread-safe
Keep your favorite Twitter threads safe by downloading a local copy

Proper README coming soon, but for now...

Download any single-author thread by idenifying the *last* tweet in the thread.
The tweet and all media files are saved and rendered as html to be displayed in your browser.

Limitations of the twitter API necessitate specifying the *last* tweet rather than the more natural first.


```
> thread-safe --help
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

```
> thread-safe thread --help
thread saves thread content and generates a local html file

Usage:
  thread-safe thread [flags] <name> <last-tweet>

Args:
  name           string  name to use for the thread
  last-tweet     string  URL or ID of the last tweet in a single-author thread

Flags:
  --no-attachments  do not download attachments

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (curre
```

```
> thread-safe html --help
html regenerates an html file from a previously saved thread

Usage:
  thread-safe html <name>

Args:
  name  string  name given to the thread

Environment Variables:
  THREAD_SAFE_PATH	top level path for thread files (current directory if unset)
  THREAD_SAFE_TOKEN	bearer token for Twitter API
```
