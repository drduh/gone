# gone

gone is an ephemeral file sharing service written in Go.

*Important* For experimental and development use only - not yet fit for production.

# Development

To build and run the application on Linux:

```
make run
```

Binaries are built into a local `release` directory for distribution and installation.

# Use

Output is structured in JSON format and can be easily parsed with `jq` for convenience, for example:

```
gone | jq .data
```

The optional `-debug` flag can be used for additional verbose program output.

# Configuration

gone uses an embedded JSON-based configuration [config/defaultSettings.json](https://github.com/drduh/gone/blob/main/config/defaultSettings.json) as default settings.

Copy the JSON file and use the `-config` flag to override options:

```
gone -config=mySettings.json
```
