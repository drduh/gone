# Design

gone is an ephemeral content sharing server written in Go.

The primary goal is to enable sharing of a link, document or other file to another device on the network.

## Features

- In-memory storage - no file artifacts after program exits
- Share files with upload, download and list functionality
- Files expire (removed) after a number of downloads or duration of time
- JSON-based configuration, logging and server responses
- Token (string-based) authentication
- Request rate and size limits
- Basic user interface without scripts
- Share and clear short text messages
- No third-party dependencies

## Security

gone has not yet been subject to security audit and has several obvious vulnerabilities, such as lack of input validation.

Therefore, gone is only intended for use on a secure network with trusted devices, such as a private LAN.

# Development

gone requires [Go 1.24.1](https://go.dev/doc/install) to develop.

## Build

To build the application on Linux:

```
make build
```

Binaries are built to `release` for distribution and installation.

## Run

To run the application on Linux:

```
make run
```

## Debug

To run the application on Linux in debug mode:

```
make debug
```

# Output

Application output is structured in JSON format and can be parsed with `jq` for convenience, for example:

```
./gone | jq .data
```

The optional `-debug` flag can be used for debug mode (additional verbose program output).

# Configuration

gone uses an embedded JSON-based configuration [config/defaultSettings.json](https://github.com/drduh/gone/blob/main/config/defaultSettings.json) as default settings.

Copy the JSON file and use the `-config` flag to set server options:

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
curl -F "file=@test.txt" http://localhost:8080/upload
```

With 3 allowed downloads before file expiration:

```
curl -F "downloads=3" -F "file=@test.txt" http://localhost:8080/upload
```

With a 15 minutes file expiration:

```
curl -F "duration=15m" -F "file=@test.txt" http://localhost:8080/upload
```

## List

List uploaded files:

```
curl http://localhost:8080/list
```

## Download

Download file (the default configuration requires basic authentication):

```
curl -H "X-Auth: mySecret" "http://localhost:8080/download/test.txt"

curl -H "X-Auth: mySecret" "http://localhost:8080/download?name=test.txt"
```

Get static (never expires) content:

```
curl http://localhost:8080/static
```

## Message

Post a message to index:

```
curl -s -F 'message=hello, world!' http://localhost:8080 >/dev/null
```

## Functions

See [config/zshrc](https://github.com/drduh/config/blob/main/zshrc#L541) for alias and function examples, such as:

```
$ gone_put test.txt 1 20m
{
  "name": "test.txt",
  "size": "6.00 Bytes",
  "owner": {
    "address": "127.0.0.1:4306",
    "agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.3"
  },
  "time": {
    "allow": "20m0s",
    "upload": "2025-03-20T10:00:00.00"
  },
  "downloads": {
    "allow": 1
  }
}
```
