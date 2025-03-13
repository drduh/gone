*Important* For experimental and development use only - not yet fit for production.

gone is an ephemeral file sharing service written in Go.

# Features

- Upload, download and list files
- File expiration (removal) after downloads or duration of time
- JSON-based configuration, logging and server responses
- Token (string-based) authentication
- Request rate-limiting
- No third-party dependencies

# Development

To build and run the application on Linux:

```
make run
```

Binaries are built into a local `release` directory for distribution and installation.

# Server

Output is structured in JSON format and can be easily parsed with `jq` for convenience, for example:

```
gone | jq .data
```

The optional `-debug` flag can be used for additional verbose program output.

## Configuration

gone uses an embedded JSON-based configuration [config/defaultSettings.json](https://github.com/drduh/gone/blob/main/config/defaultSettings.json) as default settings.

Copy the JSON file and use the `-config` flag to set server options:

```
gone -config=mySettings.json
```

# Client

The server provides a basic user interface for uploading, downloading and listing files at the default path (`/`).

All features are also available using command line programs such as curl:

Get server status:

```
curl localhost:8080/heartbeat
```

Upload file:

```
curl -F "file=@test.txt" http://localhost:8080/upload
```

Upload file with explicit number of allowed downloads before expiration:

```
curl -F "downloads=3" -F "file=@test.txt" http://localhost:8080/upload
```

List uploaded files:

```
curl http://localhost:8080/list
```

Download file (the default configuration requires basic authentication):

```
curl -H "X-Auth: mySecret" "http://localhost:8080/download?name=test.txt"
```

Get static (never expires) content:

```
curl http://localhost:8080/static
```
