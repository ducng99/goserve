# goserve - serves static files

A tool to serve static files with configurations for HTTPS, CORS and more. Inspired by PHP Dev server.

## Downloads

### Prebuilt
Prebuilt binaries can be found for Linux, Windows and MacOS - amd64/arm64/arm in [Releases](https://r.tomng.dev/goserve/releases/latest)

### Go toolchain
You can also install through go with the command below

```bash
go install r.tomng.dev/goserve@latest
```

## Usage

### Default
To start serving the current directory, you can run without providing any arguments.
It will start listening on `0.0.0.0:8080` - accepting all connections from port 8080.

```bash
goserve
```

### Custom host:port
Goserve accepts an argument as host:port to listen on.

The commands below accepts local connections on port 1337

```bash
goserve localhost:1337
```

You can also provide a host or port only, and it will fill in the default value.

Listen on `0.0.0.0:1337`

```bash
goserve :1337
```

Listen on localhost:8080

```bash
goserve localhost
```

### Help

Access `--help` anytime for info on flags allowed

### License

[MIT](./LICENSE)
