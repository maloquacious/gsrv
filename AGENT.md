# AGENT.md

## Purpose
- GSRV implements a wrapper around Go's http.Server and provides methods for shutting the server down gracefully.

## Code Style
- Standard Go formatting using `gofmt`
- Imports organized by stdlib first, then external packages
- Error handling: return errors to caller
- Function comments use Go standard format `// FunctionName does X`
- Variable naming follows camelCase
- File structure follows standard Go package conventions
