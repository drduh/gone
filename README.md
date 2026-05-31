gone is an ephemeral content server written in [Go](https://go.dev/).

The primary goal is to share files and text using an HTML interface and API.

[![test-and-lint](https://github.com/drduh/gone/actions/workflows/test-and-lint.yml/badge.svg)](https://github.com/drduh/gone/actions/workflows/test-and-lint.yml)

## Features

- No disk storage (memory only) - content cleared on exit
- No third-party dependencies - only Go required to build
- No Javascript required to use HTML-based user interface
- Share multiple files: upload, download and list feature
- Files expire after number of downloads or time duration
- Share plain-text messages and shared text area for edit
- JSON-based configurations, logging and server responses
- Token (string-based) authentication and request limiter

# Development

gone requires [Go](https://go.dev/doc/install) to develop.

[Makefile](https://github.com/drduh/gone/blob/main/Makefile) provides functionality to build, run and test the application.

## Build

```bash
make build
```

Binaries are built to the `release` directory.

## Run

```bash
make run
```

## Debug

```bash
make debug
```

## Install (Linux)

Install as a service on systemd Linux:

```bash
make install
```

# Output

Output is in JSON format and can be parsed with `jq`:

```bash
gone | jq '.message'
```

# Configuration

[Default settings](https://github.com/drduh/gone/blob/main/settings/defaultSettings.json) are embedded into the application. To configure, pass a modified settings file using `-config` or `-settings`:

```bash
gone -config mySettings.json
```

Set an empty handler path to disable it; for example, to turn off text features:

```json
"paths": {
  "message": "",
  "wall": ""
```

# Clients

A basic HTML user interface is available at the `root` path (`/` by default) - [127.0.0.1:8080](http://127.0.0.1:8080/) when running locally.

Features are also available using command-line programs such as curl:

## Status

Get server status:

```bash
curl 127.0.0.1:8080/status
```

Get user request information:

```bash
curl 127.0.0.1:8080/user
```

## Upload

Upload file:

```bash
curl 127.0.0.1:8080/upload -F "file=@example.txt"
```

Upload multiple files:

```bash
curl 127.0.0.1:8080/upload -F "file=@example1.txt" -F "file=@example2.txt"
```

Upload file with 5 downloads allowed:

```bash
curl 127.0.0.1:8080/upload -F "downloads=5" -F "file=@example.txt"
```

Upload file with 5 minute expiration:

```bash
curl 127.0.0.1:8080/upload -F "duration=5m" -F "file=@example.txt"
```

## List

List uploaded files:

```bash
curl 127.0.0.1:8080/list
```

## Download

Download a file ([default settings](https://github.com/drduh/gone/blob/main/settings/defaultSettings.json) require token-based authentication):

```bash
curl 127.0.0.1:8080/download/example.txt -H "X-Auth: mySecret"
```

Get static (never expires) content:

```bash
curl 127.0.0.1:8080/static
```

## Message

Post a plain-text message (use single quotes to wrap special characters):

```bash
curl 127.0.0.1:8080/msg -d 'message=hello, world!'
```

Get message text only:

```bash
curl 127.0.0.1:8080/msg | jq '.[].data'
```

## Wall

Post multi-line text for shared edit:

```bash
curl 127.0.0.1:8080/wall -F "wall=$(cat /etc/resolv.conf)"
```

Get shared multi-line text:

```bash
curl 127.0.0.1:8080/wall | jq -r
```

## Random

Get [random values](https://github.com/drduh/gone/blob/main/util/random.go) of certain type:

```bash
curl 127.0.0.1:8080/random/

curl 127.0.0.1:8080/random/coin

curl 127.0.0.1:8080/random/hex

curl 127.0.0.1:8080/random/id

curl 127.0.0.1:8080/random/mask

curl 127.0.0.1:8080/random/name

curl 127.0.0.1:8080/random/nato

curl 127.0.0.1:8080/random/number

curl 127.0.0.1:8080/random/pass
```

## Functions

See [config/zshrc](https://github.com/drduh/config/blob/main/zshrc#L614) for alias and function examples, such as:

```bash
$ gonePut test.txt 3 30m
[
  {
    "id": "1J81kxgMK0JEJa5VpMb7AJJvwutwaq7bhV26xtEaFL4w",
    "name": "test.txt",
    "sum": "4dca0fd5f424a31b03ab807cbae77eb32bf2d089eed1cee154b3afed458de0dc",
    "downloads": {
      "allow": 3
    },
    "size": "14 bytes",
    "type": "text/plain; charset=utf-8",
    "owner": {
      "address": "127.0.0.1:12345",
      "mask": "Bob123",
      "agent": "curl/8.7.1"
    },
    "time": {
      "allow": "30m0s",
      "upload": "2026-05-20T12:00:00.00000-00:00"
    }
  }
]
```

# Documentation

Application documentation is available with [godoc](https://go.dev/blog/godoc):

```bash
make doc
```

# Testing

Tests and lint are validated with a [workflow](https://github.com/drduh/gone/blob/main/.github/workflows/test-and-lint.yml) on changes, or manually:

```bash
make lint

make test

make test-verbose

make test-race
```

Test coverage is also available - to generate an HTML report as `testCoverage.html`:

```bash
make cover
```

# Container

The application can run in a container using [`Dockerfile`](https://github.com/drduh/gone/blob/main/Dockerfile)

On macOS, using [apple/container](https://github.com/apple/container):

```bash
make build-container

make run-container
```

Get the server IP address:

```bash
container ls | grep gone | grep --color --text -Eo "([0-9]{1,3}\.){3}[0-9]{1,3}"

curl 192.168.64.3:8080
```
