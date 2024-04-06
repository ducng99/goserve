# goserve

![goserve](https://github.com/ducng99/goserve/assets/49080794/c401087b-4130-4b52-95ab-324956e8569c)

<!-- I need a gopher here :( -->
<p align="center">
  <a href="https://github.com/ducng99/goserve/actions/workflows/build.yml">
    <img src="https://github.com/ducng99/goserve/actions/workflows/build.yml/badge.svg"/>
  </a>
  <a href="https://github.com/ducng99/goserve/actions/workflows/test.yml">
    <img src="https://github.com/ducng99/goserve/actions/workflows/test.yml/badge.svg"/>
  </a>
  <br>
  <strong>goserve</strong> is a tool to serve static files with indexing page, allows configurations for HTTPS, CORS and more.<br>Inspired by PHP Dev server.
</p>

---

## Downloads

### Prebuilt
Prebuilt binaries can be found for Linux, Windows and MacOS (amd64/arm64/arm) in [Releases](https://github.com/ducng99/goserve/releases/latest)

### Go toolchain
If you have go installed, you can also install with the command below.

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

```bash
# Listen on `0.0.0.0:1337`
goserve :1337
```

```bash
# Listen on `localhost:8080`
goserve localhost
```

```bash
# IPv6 is supported
goserve "[::0]:9876"
```

### Help

Access `--help` anytime for info on flags allowed

### License

[MIT](./LICENSE)
