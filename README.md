# tbank-installment-go-client

Go client for the T-Bank Forma installment API methods documented in this repository.

## Scope

This package implements the methods and payloads described in [doc.md](./doc.md):

- `Create`
- `Commit`
- `Cancel`
- `Info`
- webhook payload parsing helpers

The public API intentionally mirrors the documentation structure, including separate response types for different methods where the docs describe them separately.

## Installation

```bash
go get <your-module-path>
```

Before publishing, replace `<your-module-path>` with the final public module path in `go.mod`.

## Usage

```go
package main

import (
	"context"
	"log"

	tbankinstallment "tbank-installment-go-client"
)

func main() {
	client, err := tbankinstallment.NewClient(tbankinstallment.Config{
		Username:             "partner-login",
		Password:             "partner-password",
		WebhookTrustedSubnet: "91.194.226.0/23",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Create(context.Background(), tbankinstallment.CreateRequest{
		ShopID:     "SHOP_ID",
		ShowcaseID: "SHOWCASE_ID",
		Sum:        100000,
		Items: []tbankinstallment.CreateItem{
			{
				Name:     "iPhone",
				Quantity: 1,
				Price:    100000,
			},
		},
		OrderNumber: "1234567890",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.ID, resp.Link)
}
```

## Webhooks

Webhook payloads can be decoded with `ParseWebhook`. Source IP checks are optional and configured via `Config.WebhookTrustedSubnet`.

```go
payload, err := tbankinstallment.ParseWebhook(r.Body)
if err != nil {
	// handle bad payload
}

if !client.IsTrustedWebhookRequest(r) {
	// handle untrusted source
}

_ = payload
```

## Development

Run tests:

```bash
go test ./...
```

Or use the local helper targets:

```bash
make fmt
make vet
make test
make check
```

## Compatibility

- Go `1.26`
- No external runtime dependencies

## CI

GitHub Actions runs formatting checks, `go vet`, and `go test ./...` for pushes and pull requests.

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE).
