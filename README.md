# yoru

A personal library of Go binaries by ShashankRaoCoding.

This repository contains small command-line tools and utilities written in Go. Each binary is organized as a separate command under `cmd/<name>`.

## Contents

- `cmd/` — each subdirectory holds a Go `main` package for a binary.
- `pkg/` — (optional) shared packages used by the commands.

> Note: The exact layout and binaries available may change — check the `cmd/` directory for the current list of tools.

## Requirements

- Go 1.20+ (or the version listed in the repository if specified)

## Build & Install

Build a specific binary (from repository root):

```bash
# build a binary called 'example' located at cmd/example
go build -o bin/example ./cmd/example
```

Install a binary using `go install` (recommended for Go modules):

```bash
# from anywhere
go install github.com/ShashankRaoCoding/yoru/cmd/example@latest
```

Or install all binaries by building them into a local `bin/` folder:

```bash
mkdir -p bin
for d in ./cmd/*; do
  name=$(basename "$d")
  go build -o bin/$name ./cmd/$name
done
```

## Usage

Each binary typically exposes `--help` or `-h`. Example:

```bash
./bin/example --help
# or, if installed via go install
example --help
```

## Contributing

Contributions are welcome. Typical ways to contribute:

- Open an issue to discuss a bug or feature.
- Send a pull request with a clear description of the change.
- Add tests where appropriate and keep changes small and focused.

When contributing, please follow standard Go formatting and linting:

```bash
gofmt -w .
go vet ./...
```

## Reporting Issues

If you encounter a bug or unexpected behavior, please open an issue with steps to reproduce, your Go version, and any relevant logs or command output.

## Repository maintainer

ShashankRaoCoding

## License

If you want this project to be used under a specific license, add a LICENSE file in the repository. If one already exists, please refer to it for license details.

---

If you'd like, I can:

- Add a short list of the current binaries (I can scan `cmd/` and add them to the README).
- Add a LICENSE file (you can tell me which license to use).
- Create usage examples for specific binaries.
