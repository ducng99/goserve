# Codebase Overview

## Project Structure

This is a Go-based project with the following typical structure:

```
.
├── src/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal/
│       ├── pkg/
│       └── cmd/
├── test/
│   ├── test.go
│   └── files/
├── .gitignore
├── README.md
└── .continue/
    └── README.md  ← This file
```

## Key Components

### Main Application
- **src/main.go**: Entry point of the application
- **go.mod**: Go module definition and dependencies
- **go.sum**: Go module checksums

### Internal Packages
- **internal/pkg/**: Internal packages containing core logic
- **internal/cmd/**: Command-line interface implementations

### Testing
- **test/**: Test files and test data
- **test/files/**: Test data files

## Build and Development Commands

### Build
```bash
go build -v
```

### Lint
```bash
golangci-lint run
```

### Test
```bash
go test -v ./test/...
go test -v ./test/files/...
```

### Generate
```bash
go generate ./...
```

## Code Style Guidelines

- Use gofmt for formatting
- Import groups: standard library, external libraries, internal packages
- Use camelCase for variable names
- Use descriptive names, avoid single letter variables except for loop counters
- Handle errors by logging them and exiting, or returning the error
- Use explicit type declarations when type is not obvious from context
- Prefer short variable declarations (:=) over var statements
- Use consistent indentation (tabs for go files)

## Development Environment

The project is designed to run on Windows machines with Go installed. The tools and commands are compatible with Windows development environments.

## Dependencies

The project uses standard Go modules for dependency management. External dependencies are managed through go.mod and go.sum files.