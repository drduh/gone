# gone

gone is an ephemeral file sharing service written in Go.

*Important* For experimental and development use only - not yet fit for production.

# Development

To build and run the application on Linux:

```
make
```

Binaries are built into a local `release` directory for distribution and installation.

# Use

Output is structured in JSON format and can be easily parsed with `jq` for convenience, for example:

```
gone | jq .data
```

# Configuration

gone uses an embedded JSON-based configuration [config/defaultSettings.json](https://github.com/drduh/gone/blob/main/config/defaultSettings.json) as default settings.

Copy the JSON file and pass its path to gone to override configuration options, such as listening port:

```
gone -config=mySettings.json
```
