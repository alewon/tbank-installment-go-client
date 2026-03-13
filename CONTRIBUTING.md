# Contributing

Thanks for contributing.

## Guidelines

- Keep the code aligned with [doc.md](./doc.md).
- Do not merge separately documented API structures into one shared type unless the documentation does so.
- Prefer the standard library unless a dependency is clearly justified.
- Add or update tests for behavior changes.

## Local checks

```bash
gofmt -w .
go vet ./...
go test ./...
```

Or:

```bash
make check
```

## Pull requests

- Describe which documented method or structure changed.
- Mention any assumptions introduced because the upstream documentation is incomplete.
