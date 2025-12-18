# Project Analysis: GoServe - HTTP(S) Server and Proxy

## Overview

`goserve` is a versatile Go-based HTTP(S) server tool designed to quickly serve static files with directory indexing, while also supporting HTTPS, CORS, and reverse proxy functionality. It's inspired by PHP's built-in development server and provides a simple way to spin up web servers for development and testing purposes.

## Key Features

### Static File Serving
- Serve static files from any directory
- Directory listing with index pages
- Default listening on `0.0.0.0:8080`
- Support for custom host:port configurations

### HTTPS Support
- HTTPS server capability with configurable certificates
- Auto-generation of self-signed certificates
- Support for custom certificate and private key files
- Certificate storage in temporary directory
- Warning about self-signed certificates for production use

### CORS Support
- Optional CORS headers configuration
- Default headers include:
  - `Access-Control-Allow-Origin: *`
  - `Access-Control-Allow-Methods: *`
  - `Access-Control-Allow-Headers: Content-Type, Authorization`
  - `Access-Control-Max-Age: 3600`

### Reverse Proxy Functionality
- Act as a reverse proxy server
- Forward requests to target URLs
- Client IP forwarding with `X-Forwarded-For` and `X-Forwarded-Proto` headers
- Redirect handling with optional ignore flag
- Preserves original location in `X-Original-Location` header

### Customization Options
- Theme selection for directory index pages (basic or pretty with TailwindCSS)
- Log color configuration
- Custom root directory with `-d` flag
- Flexible host:port configuration

## Project Structure

Based on the README and typical Go project structure:

```
.
├── internal/
│   ├── server/          # HTTP server implementation
│   ├── proxy/           # Reverse proxy logic
│   ├── logger/          # Logging functionality
│   ├── responsewriter/  # Response handling
│   ├── ssl/             # SSL/TLS handling
│   ├── tmpl/            # Template files for directory listing
│   └── server/middlewares/  # Middleware components (e.g., CORS)
├── cmd/                 # Command-line interface
├── src/                 # Main application entry point
└── .continue/           # Documentation files
```

## Usage Examples

### Basic Usage
```bash
# Serve current directory on default port (0.0.0.0:8080)
goserve

# Serve on specific host:port
goserve localhost:1337

# Serve specific directory
goserve -d ./web/static/
```

### HTTPS
```bash
# Start HTTPS server with auto-generated certificate
goserve -s

# Start HTTPS with custom certificate
goserve -s --sslcert /path/to/cert.crt --sslkey /path/to/priv.key
```

### Proxy Mode
```bash
# Forward requests to target URL
goserve -p http://pi.local/ :9999

# Proxy with headers
goserve -p http://localhost:8080 --proxy-headers localhost:8081
```

## Installation

### Prebuilt Binaries
Available for Linux, Windows, and MacOS (amd64/arm64/arm) in [Releases](https://github.com/ducng99/goserve/releases/latest)

### Go Toolchain
```bash
go install github.com/ducng99/goserve@latest
go run github.com/ducng99/goserve@latest
```

## Target Audience

This tool is primarily designed for:
- Development environments where quick static file serving is needed
- Testing web applications locally
- Creating secure development servers with HTTPS
- Reverse proxying for local development setups
- Simple file sharing between team members

## Build and Test

```bash
# Build
go build -v

# Test
go test -v ./test/...

# Lint
golangci-lint run
```

## License

MIT License