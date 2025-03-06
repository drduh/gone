# gone

gone is an ephemeral file sharing service written in Go.

*Important* For experimental and development use only - not yet fit for production.

# Development

To build and run the application on Linux:

```
make
```

Output is structured in JSON format and can be easily parsed with `jq` for convenience, for example:

```
make | jq .data
```
