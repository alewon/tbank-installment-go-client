GO ?= go

.PHONY: test fmt vet check

test:
	$(GO) test ./...

fmt:
	gofmt -w .

vet:
	$(GO) vet ./...

check: fmt vet test
