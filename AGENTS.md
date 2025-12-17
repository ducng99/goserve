# AGENTS.md

Build/lint/test commands:
- Build: `go build -v`
- Lint: `golangci-lint run` (if installed)
- Test: `go test -v ./test/...` or `go test -v ./test/files/...` for specific tests
- Generate: `go generate ./...`

Code style guidelines:
- Use gofmt for formatting
- Import groups: standard library, external libraries, internal packages
- Use camelCase for variable names
- Use descriptive names, avoid single letter variables except for loop counters
- Handle errors by logging them and exiting, or returning the error
- Use explicit type declarations when type is not obvious from context
- Prefer short variable declarations (:=) over var statements
- Use consistent indentation (tabs for go files)