# Go-клиент для сервиса «T-Рассрочка»

## Установка

```bash
go get github.com/alewon/tbank-installment-go-client
```

## Быстрый старт

```go
package main

import (
	"context"
	"log"

	tbankinstallment "github.com/alewon/tbank-installment-go-client"
)

func main() {
	client, err := tbankinstallment.NewClient(tbankinstallment.Config{
		Username: "showcase-id",
		Password: "password",
		Demo:     true,
		BaseURL:  "https://forma.tbank.ru",
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.CreateDemo(context.Background(), tbankinstallment.CreateDemoRequest{
		CreateRequestBody: tbankinstallment.CreateRequestBody{
			ShopID:      "SHOP_ID",
			ShowcaseID:  "SHOWCASE_ID",
			Sum:         3200,
			OrderNumber: "order-123",
			PromoCode:   "installment_0_0_3_5,41",
			Items: []tbankinstallment.CreateItem{
				{
					Name:     "Товар",
					Quantity: 1,
					Price:    3200,
				},
			},
		},
		DemoFlow: tbankinstallment.DemoFlowSMS,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.ID, resp.Link)
}
```

Для получения информации по тестовой заявке:

```go
info, err := client.Info(context.Background(), "order-123")
if err != nil {
	log.Fatal(err)
}

log.Println(info.Status, info.Committed)
```
