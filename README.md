# goserve

![goserve](https://github.com/ducng99/goserve/assets/49080794/8973ff0c-1b73-4d0d-8d47-792594ca8005)

<!-- I need a gopher here :( -->
<p align="center">
  <a href="https://github.com/ducng99/goserve/actions/workflows/build.yml">
    <img src="https://github.com/ducng99/goserve/actions/workflows/build.yml/badge.svg"/></a>
  <a href="https://github.com/ducng99/goserve/actions/workflows/test.yml">
    <img src="https://github.com/ducng99/goserve/actions/workflows/test.yml/badge.svg"/>
  </a>
  <br>
  <strong>goserve</strong> is a quick tool to spin up a HTTP(S) server to host static files with directory indexing page, with configurable HTTPS, CORS and more.<br>Inspired by PHP Dev server.
</p>

---

## Installation

### Prebuilt
Prebuilt binaries can be found for Linux, Windows and MacOS (amd64/arm64/arm) in [Releases](https://github.com/ducng99/goserve/releases/latest)

### Go toolchain
If you have go installed, you can also install with the command below.

```bash
go install r.tomng.dev/goserve@latest
```

Or just run it without installing

```bash
go run r.tomng.dev/goserve@latest
```

## Usage & examples

### Default
To start serving the current directory, you can run without providing any arguments.
It will start listening on `0.0.0.0:8080` - accepting all connections from port 8080.

```bash
goserve
```

### Custom host:port
goserve accepts an argument as host:port to listen on.

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
goserve "[::0]:1337"
```

### Different root dir

By default, goserve uses the current directory and serve its files and directories. You can change it using `-d` or `--dir` flag.

```bash
goserve -d ./web/static/
```

### HTTPS and certificates
goserve can start a HTTPS server with your provided certificate and private key, or generate a pair if you don't.
Generated certificate and key is stored in `[TempDir]/goserve/` directory.

**Note:** If you want to bypass the warning in browsers when accessing self-signed site, it is recommended to use tools like [mkcert](https://github.com/FiloSottile/mkcert) to set up local CA.
You can then pass its certificate and private key to goserve.

> [!WARNING]
> You should not use self-signed certificate in production environment, it should only be used for local development testing.
> The private key should not be shared.

```bash
# Starts HTTPS server and use auto-generated self-signed certificate and key
goserve -s
```

```bash
# Starts HTTPS server with provided certificate and key
goserve -s --sslcert /path/to/cert.crt --sslkey /path/to/priv.key
```

### CORS

CORS headers aren't added by default when serving files, you can supply `--cors` flag to add these headers.

```bash
goserve --cors
```

[Headers](./internal/server/middlewares/cors.go):
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: *
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 3600
```

### Proxy

Instead of serving your local files and showing directory listing, goserve can act as a reverse proxy server and forward requests to a target URL.
This is useful if you want to use goserve HTTPS on top of your HTTP web server.

Use `-p` or `--proxy` with the target URL you want to forward to.

```bash
# Listen on port 9999 and forward requests to http://pi.local
goserve -p http://pi.local/ :9999
```

#### Client IP forwarding
If `--proxy-headers` flag is set, goserve includes `X-Forwarded-For` and `X-Forwarded-Proto` headers, setting client's IP and the original protocol used, respectively.

#### Redirect response
Sometimes target server can return a redirect through `Location` header.
goserve can strips out this header if `--proxy-ignore-redirect` flag is specified.

The original `Location` header value can be found in `X-Original-Location` header.

### Theme

#### Directory index page
By default, directory index page uses "pretty" theme with TailwindCSS. You can switch to "basic" theme by suppling `--index-theme` flag, which contains just a simple HTML page (with very minimal CSS).

#### Log color
If you prefer default text color only for logs, setting `--log-color=false` flag will disable all colors when logging.

## Help

Access `--help` anytime for up-to-date info on flags

```bash
$ goserve -h
Starts a web server to serve static files, with options for HTTPS, directory, CORS, and more.

Usage:
  goserve [flags] [host:port]

Default host:port is "0.0.0.0:8080"

Examples:
Start server with HTTPS on port 8443:
goserve -cd /path/to/dir --https --sslcert full-cert.crt --sslkey private-key.key localhost:8443

Proxy to another server on port 8080, and listen on port 8081:
goserve -p http://localhost:8080 localhost:8081

Flags:
  -c, --cors                    Set CORS headers
  -d, --dir string              Directory to serve (default ".")
  -h, --help                    help for goserve
      --https                   Alias for --ssl
      --index-theme string      Directory index page theme.
                                Available themes: basic, pretty (default "pretty")
      --log-color               Disable colored log output (default true)
  -p, --proxy string            Proxy forward to the specified URL.
                                This will disable directory listing and file serving.
      --proxy-headers           Include X-Forwarded-For and X-Forwarded-Proto headers in proxy request (default true)
      --proxy-ignore-redirect   Ignore redirects from the target server
  -s, --ssl                     Use HTTPS server
      --sslcert string          Path to a full certificate file
      --sslkey string           Path to a private key file
  -v, --version                 version for goserve
```

## License

[MIT](./LICENSE)
