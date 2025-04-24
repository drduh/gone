# Design

gone is an ephemeral content hosting server written in [Go](https://go.dev/).

The primary goal is to enable sharing of files and text using command-line API or simple HTML user interface.

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

## Build

To build on Linux:

```
make build
```

Binaries are built to the `release` directory for distribution and installation.

## Run

To run on Linux:

```
make run
```

## Debug

To run in debug mode on Linux:

```
make debug
```

# Output

Application output is structured in JSON format and can be parsed with `jq` for convenience, for example:

```
./gone | jq .data
```

The optional `-debug` flag can be used for debug mode (provides additional application output).

# Configuration

gone uses an embedded JSON-based configuration [defaultSettings.json](https://github.com/drduh/gone/blob/main/settings/defaultSettings.json) as default settings.

To change application settings, copy the default settings JSON file and use the `-config` flag:

```
./gone -config=mySettings.json
```

# Client

The server provides a basic user interface for uploading, downloading and listing files at the default path (`/`):

[localhost:8080](http://localhost:8080)

All features are also available using command line programs such as curl:

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

Download file (the [default configuration](https://github.com/drduh/gone/blob/main/config/defaultSettings.json) requires basic authentication):

```
curl -H "X-Auth: mySecret" "localhost:8080/download/test.txt"
```

Get static (never expires) content:

```
curl localhost:8080/static
```

## Message

Post a message (use single quotes to wrap special characters):

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

curl localhost:8080/random/name

curl localhost:8080/random/nato

curl localhost:8080/random/number

curl localhost:8080/random/coin
```

## Functions

See [config/zshrc](https://github.com/drduh/config/blob/main/zshrc#L541) for alias and function examples, such as:

```
$ gone_put test.txt 1 20m
[
  {
    "name": "test.txt",
    "downloads": {
      "allow": 1
    },
    "size": "6.00 Bytes",
    "owner": {
      "address": "127.0.0.1:4306",
      "agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.3"
    },
    "time": {
      "allow": "20m0s",
      "upload": "2025-04-22T10:00:00"
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

Unit tests are validated by a [workflow](https://github.com/drduh/gone/blob/main/.github/workflows/test-and-lint.yml) on repository changes. They can also be run manually:

```
make test

make test-verbose
```

Test coverage is also available - to generate an HTML report (`testCoverage.html`):

```
make cover
```
