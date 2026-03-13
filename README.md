# Go-клиент для сервиса «T-Рассрочка»

## Что реализовано

Пакет реализует методы и структуры данных, описанные в [doc.md](./doc.md):

- `Create`
- `Commit`
- `Cancel`
- `Info`
- вспомогательные функции для разбора webhook payload

Публичный API намеренно повторяет структуру документации, включая отдельные типы ответов для разных методов там, где это разделено в исходном описании.

## Установка

```bash
go get github.com/alewon/tbank-installment-go-client
```

## Использование

```go
package main

import (
	"context"
	"log"

	tbankinstallment "github.com/alewon/tbank-installment-go-client"
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

## Webhook

Webhook payload можно декодировать через `ParseWebhook`. Проверка IP-адреса источника опциональна и настраивается через `Config.WebhookTrustedSubnet`.

```go
payload, err := tbankinstallment.ParseWebhook(r.Body)
if err != nil {
	// обработать некорректный payload
}

if !client.IsTrustedWebhookRequest(r) {
	// обработать недоверенный источник
}

_ = payload
```

## Разработка

Для запуска тестов:

```bash
go test ./...
```

## Совместимость

- Go `1.18+`
- внешние runtime-зависимости отсутствуют

## CI

GitHub Actions запускает проверки форматирования, `go vet` и `go test ./...` для `push` и `pull request`.

## Лицензия

Проект распространяется по лицензии MIT. Подробности см. в [LICENSE](./LICENSE).
