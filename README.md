gone is an ephemeral content hosting server written in [Go](https://go.dev/).

The primary goal is to share files or text using an HTML interface or an API.

## Features

- No disk storage (memory only) - content cleared on exit
- No third-party dependencies - only Go required to build
- No Javascript required to use HTML-based user interface
- Share multiple files: upload, download and list feature
- Files expire after number of downloads or time duration
- Share short text messages and shared text area for edit
- JSON-based configurations, logging and server responses
- Token (string-based) authentication and request limiter

# Development

gone requires [Go](https://go.dev/doc/install) to develop.

[Makefile](https://github.com/drduh/gone/blob/main/Makefile) provides functionality to build, run and test the application.

## Build

```
make build
```

Binaries are built to the `release` directory for distribution and installation.

## Run

```
make run
```

## Debug

```
make debug
```

## Install

To install as a service on systemd Linux:

```
make install
```

# Output

Output is in JSON format and can be parsed with `jq`:

```
gone | jq '.message'
```

# Settings

[Default settings](https://github.com/drduh/gone/blob/main/settings/defaultSettings.json) are embedded into the application. To make changes, copy and pass a modified settings file using `-config`:

```
gone -config mySettings.json
```

Set an empty handler path to disable it; for example, to turn off text features:

```
"paths": {
  "message": "",
  "wall": ""
```

# Clients

A basic HTML user interface is available at the `root` path (`/` by default) - [localhost:8080](http://localhost:8080/) when running locally.

Features are also available using command-line programs such as curl:

## Status

Get server status:

```
curl localhost:8080/heartbeat
```

## Upload

Upload file:

```
curl -F "file=@test.txt" localhost:8080/upload
```

Upload multiple files:

```
curl -F "file=@test.txt" -F "file=@test2.txt" localhost:8080/upload
```

With 3 allowed downloads before file expiration:

```
curl -F "downloads=3" -F "file=@test.txt" localhost:8080/upload
```

With a 15 minutes file expiration:

```
curl -F "duration=15m" -F "file=@test.txt" localhost:8080/upload
```

## List

List uploaded files:

```
curl localhost:8080/list
```

## Download

Download a file ([default settings](https://github.com/drduh/gone/blob/main/settings/defaultSettings.json) require token-based authentication):

```
curl -H "X-Auth: mySecret" "localhost:8080/download/test.txt"
```

Get static (never expires) content:

```
curl localhost:8080/static
```

## Message

Post a short text message (use single quotes to wrap special characters):

```
curl -s -F 'message=hello, world!' localhost:8080/msg >/dev/null
```

## Wall

Post multi-line text for shared edit:

```
curl -s -F "wall=$(cat /etc/dnsmasq.conf)" localhost:8080/wall >/dev/null
```

Get shared multi-line text:

```
curl localhost:8080/wall | jq -r
```

## Random

Get a [random value](https://github.com/drduh/gone/blob/main/util/random.go) of certain type:

```
curl localhost:8080/random/

curl localhost:8080/random/coin

curl localhost:8080/random/name

curl localhost:8080/random/nato

curl localhost:8080/random/number
```

## Functions

See [config/zshrc](https://github.com/drduh/config/blob/main/zshrc#L614) for alias and function examples, such as:

```
$ gonePut test.txt 3 30m
[
  {
    "name": "test.txt",
    "downloads": {
      "allow": 3
    },
    "size": "13.00 Bytes",
    "owner": {
      "address": "127.0.0.1:1234",
      "mask": "Bob123",
      "agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
    },
    "time": {
      "allow": "30m0s",
      "upload": "2025-10-31T10:00:00"
    }
  }
]
```

# Documentation

Application documentation is available with [godoc](https://go.dev/blog/godoc):

```
make doc
```

# Testing

Tests and lint are validated with a [workflow](https://github.com/drduh/gone/blob/main/.github/workflows/test-and-lint.yml) on changes, or manually:

```
make lint

make test

make test-verbose

make test-race
```

Test coverage is also available - to generate an HTML report as `testCoverage.html`:

```
make cover
```
