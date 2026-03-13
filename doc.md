# T-Bank Forma API: полная фиксация методов заявки и callback

Документ собран только по двум страницам:

- <https://forma.tbank.ru/docs/online/setup-types/api/methods>
- <https://forma.tbank.ru/docs/online/settings/app-notifications>

Цель этого файла: зафиксировать абсолютно все методы, структуры, поля, enum-значения, статусы, примеры и оговорки, которые явно присутствуют на исходных страницах. Если исходная документация что-то не описывает формально, это тоже зафиксировано явно, чтобы ничего не потерялось и ничего не было додумано.

## 1. Полный список методов и callback

На указанных страницах описаны:

- `Create` — создание заявки
- `Commit` — подтверждение заявки
- `Cancel` — отмена заявки
- `Info` — получение текущего состояния заявки
- `Callback / HTTP-нотификации` — уведомления об изменениях состояния заявки

Других API-методов на этих двух страницах нет.

## 2. Общие ограничения и HTTP-статусы

### 2.1. Лимит запросов

- Можно отправлять не более `25` запросов в секунду.
- При превышении лимита возвращается:
  - HTTP-статус `429`
  - заголовок `Retry-After` с количеством секунд ожидания

### 2.2. Общие HTTP-статусы для методов заявки

На странице методов явно указаны такие статусы:

- `200` — успешный сценарий
- `400` — некорректный формат запроса
- `401` — аутентификация не пройдена: введены неверные логин и/или пароль
- `403` — авторизация не пройдена: попытка работать с чужой заявкой
- `422` — ошибка бизнес-логики: в текущем состоянии заявки нельзя выполнить это действие
- `500` — внутренняя ошибка сервера

Дополнительно выше на той же странице для лимита запросов указан:

- `429` — превышен лимит запросов

## 3. Общие enum-значения, встречающиеся в ответах

### 3.1. Общий статус заявки

Эти значения статуса прямо перечислены на странице для `Commit`, `Cancel` и `Info`:

- `new` — заявка создана, покупатель еще не подтвердил свои данные кодом из СМС
- `inprogress` — покупатель подтвердил свои данные вводом кода из СМС (ФИО, номер телефона, e-mail) и находится на шаге ввода данных паспорта либо ожидает одобрения заявки хотя бы от одного банка
- `approved` — заявка одобрена хотя бы одним банком
- `signed` — заявка подписана покупателем при помощи СМС-кода, Self id или на встрече
- `canceled` — заявка отменена покупателем
- `rejected` — по заявке пришел отказ от всех банков

### 3.2. Общий тип продукта

Эти значения явно перечислены в описании поля `product`:

- `credit` — кредит
- `installment_credit` — рассрочка

### 3.3. Общие способы подписания

Эти значения явно перечислены в описании поля `appropriate_signing_types`:

- `bank` — подписание на встрече
- `sms` — подписание через СМС для повторных клиентов банка
- `ses` — подписание по Self id

## 4. Method: Create

### 4.1. Назначение

`Create` используется для создания реальных заявок.

Если метод случайно вызван повторно с одинаковым составом заказа, включая номер заказа в системе партнера `orderNumber`, новая заявка не создается. Возвращается ссылка на уже созданную заявку.

### 4.2. HTTP

- Метод: `POST`
- URL: `https://forma.tbank.ru/api/partners/v2/orders/create`

### 4.3. Важные замечания

- `orderNumber` обязателен, если используются вебхуки и методы API `Commit`, `Cancel` и `Info`.
- `promoCode` по умолчанию: `default`.
- В примере запроса на странице есть поле `demoFlow`, но в таблице параметров `Create` такого поля нет.

### 4.4. Структуры запроса

#### 4.4.1. `CreateRequestBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `shopId` | `string(50)` | required | Идентификатор компании |
| `showcaseId` | `string(50)` | required | Идентификатор магазина, сайта |
| `sum` | `number` | required | Общая сумма заказа |
| `items` | `CreateRequestBodyItem[]` | required | Список товаров |
| `orderNumber` | `string(64)` | optional | Номер заказа в системе партнера. Обязательный, если используются вебхуки и методы API `Commit`, `Cancel` и `Info` |
| `promoCode` | `string(64)` | optional | Идентификатор продукта: кредит или рассрочка |
| `webhookURL` | `string` | optional | Ссылка для отправки вебхуков |
| `successURL` | `string` | optional | Ссылка для возврата в случае успешного завершения заявки, то есть подписания |
| `failURL` | `string` | optional | Ссылка для возврата в случае неуспешного завершения заявки, то есть отказа банка |
| `returnURL` | `string` | optional | Ссылка для возврата в случае отмены заявки |
| `values` | `CreateRequestBodyValues` | optional | Данные покупателя для предзаполнения формы |

#### 4.4.2. `CreateRequestBodyItem`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `name` | `string(255)` | required | Название товарной позиции |
| `quantity` | `integer` | required | Количество единиц товара указанной позиции |
| `price` | `number` | required | Цена одной единицы товара |
| `category` | `string(255)` | optional | Категория товара |
| `vendorCode` | `string(64)` | optional | Артикул |

#### 4.4.3. `CreateRequestBodyValues`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `contact` | `CreateRequestBodyContact` | optional | Данные покупателя |

#### 4.4.4. `CreateRequestBodyContact`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `fio` | `CreateRequestBodyFio` | optional | ФИО покупателя |
| `mobilePhone` | `string` | optional | Телефон покупателя. Подходит строка с номером телефона РФ в любом формате |
| `email` | `string` | optional | E-mail покупателя |

#### 4.4.5. `CreateRequestBodyFio`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `lastName` | `string` | optional | Фамилия покупателя |
| `firstName` | `string` | optional | Имя покупателя |
| `middleName` | `string` | optional | Отчество покупателя |

### 4.5. Структуры ответа

#### 4.5.1. `CreateResponseBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `id` | `string` | required | ID заявки в TCB |
| `link` | `string` | required | Ссылка на заявку в TCB |

### 4.6. Статусы для метода Create

- `200` — заявка успешно создана, либо возвращена ссылка на уже существующую заявку при повторном вызове с тем же составом заказа и тем же `orderNumber`
- `400` — некорректный формат запроса
- `401` — неверные логин и/или пароль Basic Auth
- `429` — превышен лимит запросов, в заголовке есть `Retry-After`
- `500` — внутренняя ошибка сервера

### 4.7. Ошибки для метода Create

На странице нет структуры тела ошибки для `Create`.

### 4.8. Пример запроса Create

```bash
curl -v -XPOST -H 'Content-type: application/json' -d '{
  "shopId": "SHOP_ID",
  "showcaseId": "SHOWCASE_ID",
  "sum": 100000,
  "items": [
    {
      "name": "iPhone",
      "quantity": 1,
      "price": 100000,
      "category": "mobile",
      "vendorCode": "12345"
    }
  ],
  "orderNumber": "1234567890",
  "promoCode": "default",
  "demoFlow": "sms",
  "values": {
    "contact": {
      "fio": {
        "lastName": "Иванов",
        "firstName": "Иван",
        "middleName": "Иванович"
      },
      "mobilePhone": "9998887766",
      "email": "ivan@example.com"
    }
  }
}' 'https://forma.tbank.ru/api/partners/v2/orders/create'
```

## 5. Method: Commit

### 5.1. Назначение

`Commit` подтверждает актуальность заявки после подписания документов покупателем. После подтверждения заявки деньги придут на счет магазина.

### 5.2. HTTP

- Метод: `POST`
- URL: `https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/commit`

### 5.3. Важные замечания

- Метод должен быть вызван в течение `14` дней после подписания покупателем заявки.
- Если подключено автоподтверждение заказов, вызывать метод не нужно, заявка подтвердится автоматически.
- Если после подписания документов заявка охлаждается, вызвать `Commit` можно только после завершения периода охлаждения.
- До окончания периода охлаждения будет возвращаться ошибка `422`.

### 5.4. Структуры запроса

#### 5.4.1. `CommitRequestPath`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `orderNumber` | `string` | required | Номер заказа |

#### 5.4.2. `CommitRequestHeaders`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `Authorization` | `header` | required | Basic Auth |

#### 5.4.3. `CommitRequestBody`

Тело запроса на странице не описано.

### 5.5. Структуры ответа

#### 5.5.1. `CommitResponseBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `id` | `string(64)` | required | Номер заказа. Соответствует `orderNumber`, если передается при создании |
| `status` | `enum` | required | Статус заявки |
| `created_at` | `string dateTime` | required | Дата и время создания заявки |
| `demo` | `boolean` | required | Флаг, является ли заявка тестовой |
| `committed` | `boolean` | required | Флаг, является ли заявка подтвержденной |
| `first_payment` | `number double` | not specified | Первоначальный взнос |
| `order_amount` | `number double` | required | Сумма заказа |
| `credit_amount` | `number double` | not specified | Сумма выдаваемого покупателю кредита |
| `transfer_amount` | `number double` | not specified | Сумма для перевода партнеру |
| `product` | `enum` | not specified | Тип продукта: `credit` или `installment_credit` |
| `term` | `integer int32` | not specified | Срок кредита в месяцах |
| `monthly_payment` | `number double` | not specified | Ежемесячный платеж |
| `first_name` | `string` | not specified | Имя покупателя |
| `last_name` | `string` | not specified | Фамилия покупателя |
| `middle_name` | `string` | not specified | Отчество покупателя |
| `phone` | `string` | not specified | Телефон покупателя |
| `loan_number` | `string` | not specified | Номер кредитного договора |
| `email` | `string` | not specified | E-mail покупателя |
| `appropriate_signing_types` | `string[]` | required | Доступные покупателю способы подписания |
| `signing_type` | `string` | not specified | Актуальный способ подписания заявки. Входит в `appropriate_signing_types` |
| `chosen_bank` | `string` | not specified | Банк, выдавший кредит. ОТП Банк, Т-Банк или МТС Банк |
| `expected_overdue_at` | `string dateTime` | required | Дата и время окончания срока действия заявки |
| `items` | `CommitResponseBodyItem[]` | not specified | Состав заказа покупателя |
| `commit_cooldown` | `CommitResponseBodyCommitCooldown` | not specified | Информация о периоде охлаждения |

#### 5.5.2. `CommitResponseBodyStatusEnum`

- `new`
- `inprogress`
- `approved`
- `signed`
- `canceled`
- `rejected`

#### 5.5.3. `CommitResponseBodyProductEnum`

- `credit`
- `installment_credit`

#### 5.5.4. `CommitResponseBodyAppropriateSigningTypesEnum`

- `bank`
- `sms`
- `ses`

#### 5.5.5. `CommitResponseBodyItem`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `name` | `string(255)` | required | Название товарной позиции |
| `quantity` | `integer` | required | Количество единиц товара указанной позиции |
| `price` | `number` | required | Цена одной единицы товара |
| `category` | `string(255)` | optional | Категория товара |
| `vendorCode` | `string(64)` | optional | Артикул |

#### 5.5.6. `CommitResponseBodyCommitCooldown`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `until` | `string dateTime` | required | Дата и время окончания периода охлаждения в ISO 8601 формате |

### 5.6. Статусы для метода Commit

- `200` — успешное подтверждение заявки
- `400` — некорректный формат запроса
- `401` — неверные логин и/или пароль Basic Auth
- `403` — попытка работать с чужой заявкой
- `422` — в текущем состоянии заявки нельзя выполнить это действие
- `422` — также возвращается, если заявка еще находится в периоде охлаждения
- `429` — превышен лимит запросов, в заголовке есть `Retry-After`
- `500` — внутренняя ошибка сервера

### 5.7. Ошибки для метода Commit

На странице нет структуры тела ошибки для `Commit`.

### 5.8. Пример запроса Commit

```bash
curl -v -XPOST -H 'Authorization: Basic <TOKEN>' -H 'Content-type: application/json' 'https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/commit'
```

## 6. Method: Cancel

### 6.1. Назначение

`Cancel` отменяет заявку, если она стала неактуальной.

На странице приведены примеры:

- покупатель уже подписал кредитный договор, а товара или услуги нет в наличии
- покупатель сам обратился в магазин и попросил отменить заявку

### 6.2. HTTP

- Метод: `POST`
- URL: `https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/cancel`

### 6.3. Важные замечания

- Если заявка отменена, восстановить ее невозможно, нужно заполнять новую.
- Если подключено автоподтверждение заказов, после подписания документов покупателем отменить заявку уже нельзя.
- В таком случае придется оформлять возврат.

### 6.4. Структуры запроса

#### 6.4.1. `CancelRequestPath`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `orderNumber` | `string` | required | Номер заказа |

#### 6.4.2. `CancelRequestHeaders`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `Authorization` | `header` | required | Basic Auth |

#### 6.4.3. `CancelRequestBody`

Тело запроса на странице не описано.

### 6.5. Структуры ответа

#### 6.5.1. `CancelResponseBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `id` | `string(64)` | required | Номер заказа. Соответствует `orderNumber`, если передается при создании |
| `status` | `enum` | required | Статус заявки |
| `created_at` | `string dateTime` | required | Дата и время создания заявки |
| `demo` | `boolean` | required | Флаг, является ли заявка тестовой |
| `committed` | `boolean` | required | Флаг, является ли заявка подтвержденной |
| `first_payment` | `number double` | not specified | Первоначальный взнос |
| `order_amount` | `number double` | required | Сумма заказа |
| `credit_amount` | `number double` | not specified | Сумма выдаваемого покупателю кредита |
| `transfer_amount` | `number double` | not specified | Сумма для перевода партнеру |
| `product` | `enum` | not specified | Тип продукта: `credit` или `installment_credit` |
| `term` | `integer int32` | not specified | Срок кредита в месяцах |
| `monthly_payment` | `number double` | not specified | Ежемесячный платеж |
| `first_name` | `string` | not specified | Имя покупателя |
| `last_name` | `string` | not specified | Фамилия покупателя |
| `middle_name` | `string` | not specified | Отчество покупателя |
| `phone` | `string` | not specified | Телефон покупателя |
| `loan_number` | `string` | not specified | Номер кредитного договора |
| `email` | `string` | not specified | E-mail покупателя |
| `appropriate_signing_types` | `string[]` | required | Доступные покупателю способы подписания |
| `signing_type` | `string` | not specified | Актуальный способ подписания заявки. Входит в `appropriate_signing_types` |
| `chosen_bank` | `string` | not specified | Банк, выдавший кредит. ОТП Банк, Т-Банк или МТС Банк |
| `expected_overdue_at` | `string dateTime` | required | Дата и время окончания срока действия заявки |
| `items` | `CancelResponseBodyItem[]` | not specified | Состав заказа покупателя |
| `commit_cooldown` | `CancelResponseBodyCommitCooldown` | not specified | Информация о периоде охлаждения |

#### 6.5.2. `CancelResponseBodyStatusEnum`

- `new`
- `inprogress`
- `approved`
- `signed`
- `canceled`
- `rejected`

#### 6.5.3. `CancelResponseBodyProductEnum`

- `credit`
- `installment_credit`

#### 6.5.4. `CancelResponseBodyAppropriateSigningTypesEnum`

- `bank`
- `sms`
- `ses`

#### 6.5.5. `CancelResponseBodyItem`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `name` | `string(255)` | required | Название товарной позиции |
| `quantity` | `integer` | required | Количество единиц товара указанной позиции |
| `price` | `number` | required | Цена одной единицы товара |
| `category` | `string(255)` | optional | Категория товара |
| `vendorCode` | `string(64)` | optional | Артикул |

#### 6.5.6. `CancelResponseBodyCommitCooldown`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `until` | `string dateTime` | required | Дата и время окончания периода охлаждения в ISO 8601 формате |

### 6.6. Статусы для метода Cancel

- `200` — успешная отмена заявки
- `400` — некорректный формат запроса
- `401` — неверные логин и/или пароль Basic Auth
- `403` — попытка работать с чужой заявкой
- `422` — в текущем состоянии заявки нельзя выполнить это действие
- `429` — превышен лимит запросов, в заголовке есть `Retry-After`
- `500` — внутренняя ошибка сервера

### 6.7. Ошибки для метода Cancel

На странице нет структуры тела ошибки для `Cancel`.

### 6.8. Пример запроса Cancel

```bash
curl -v -XPOST -H 'Authorization: Basic <TOKEN>' -H 'Content-type: application/json' 'https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/cancel'
```

## 7. Method: Info

### 7.1. Назначение

`Info` возвращает актуальный статус заявки.

На странице отдельно рекомендовано:

- использовать метод отдельно или в комбинации с HTTP-нотификациями
- всегда проверять актуальный статус заказа методом `Info`
- делать это после получения вебхуков
- делать это, если вебхук не пришел в течение часа после создания заявки

### 7.2. HTTP

- Метод: `GET`
- URL: `https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/info`

### 7.3. Структуры запроса

#### 7.3.1. `InfoRequestPath`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `orderNumber` | `string` | required | Номер заказа |

#### 7.3.2. `InfoRequestHeaders`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `Authorization` | `header` | required | Basic Auth |

#### 7.3.3. `InfoRequestBody`

Тело запроса на странице не описано.

### 7.4. Структуры ответа

#### 7.4.1. `InfoResponseBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `id` | `string(64)` | required | Номер заказа. Соответствует `orderNumber`, если передается при создании |
| `status` | `enum` | required | Статус заявки |
| `created_at` | `string dateTime` | required | Дата и время создания заявки |
| `demo` | `boolean` | required | Флаг, является ли заявка тестовой |
| `committed` | `boolean` | required | Флаг, является ли заявка подтвержденной |
| `first_payment` | `number double` | not specified | Первоначальный взнос |
| `order_amount` | `number double` | required | Сумма заказа |
| `credit_amount` | `number double` | not specified | Сумма выдаваемого покупателю кредита |
| `transfer_amount` | `number double` | not specified | Сумма для перевода партнеру |
| `product` | `enum` | not specified | Тип продукта: `credit` или `installment_credit` |
| `term` | `integer int32` | not specified | Срок кредита в месяцах |
| `monthly_payment` | `number double` | not specified | Ежемесячный платеж |
| `first_name` | `string` | not specified | Имя покупателя |
| `last_name` | `string` | not specified | Фамилия покупателя |
| `middle_name` | `string` | not specified | Отчество покупателя |
| `phone` | `string` | not specified | Телефон покупателя |
| `loan_number` | `string` | not specified | Номер кредитного договора |
| `email` | `string` | not specified | E-mail покупателя |
| `appropriate_signing_types` | `string[]` | required | Доступные покупателю способы подписания |
| `signing_type` | `string` | not specified | Актуальный способ подписания заявки. Входит в `appropriate_signing_types` |
| `chosen_bank` | `string` | not specified | Банк, выдавший кредит. ОТП Банк, Т-Банк или МТС Банк |
| `expected_overdue_at` | `string dateTime` | required | Дата и время окончания срока действия заявки |
| `items` | `InfoResponseBodyItem[]` | not specified | Состав заказа покупателя |
| `commit_cooldown` | `InfoResponseBodyCommitCooldown` | not specified | Информация о периоде охлаждения |

#### 7.4.2. `InfoResponseBodyStatusEnum`

- `new`
- `inprogress`
- `approved`
- `signed`
- `canceled`
- `rejected`

#### 7.4.3. `InfoResponseBodyProductEnum`

- `credit`
- `installment_credit`

#### 7.4.4. `InfoResponseBodyAppropriateSigningTypesEnum`

- `bank`
- `sms`
- `ses`

#### 7.4.5. `InfoResponseBodyItem`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `name` | `string(255)` | required | Название товарной позиции |
| `quantity` | `integer` | required | Количество единиц товара указанной позиции |
| `price` | `number` | required | Цена одной единицы товара |
| `category` | `string(255)` | optional | Категория товара |
| `vendorCode` | `string(64)` | optional | Артикул |

#### 7.4.6. `InfoResponseBodyCommitCooldown`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `until` | `string dateTime` | required | Дата и время окончания периода охлаждения в ISO 8601 формате |

### 7.5. Статусы для метода Info

- `200` — информация по заявке получена успешно
- `400` — некорректный формат запроса
- `401` — неверные логин и/или пароль Basic Auth
- `403` — попытка работать с чужой заявкой
- `422` — в текущем состоянии заявки нельзя выполнить это действие
- `429` — превышен лимит запросов, в заголовке есть `Retry-After`
- `500` — внутренняя ошибка сервера

### 7.6. Ошибки для метода Info

На странице нет структуры тела ошибки для `Info`.

### 7.7. Пример запроса Info

```bash
curl -v -H 'Authorization: Basic <TOKEN>' -H 'Content-type: application/json' 'https://forma.tbank.ru/api/partners/v2/orders/{orderNumber}/info'
```

## 8. Callback / HTTP-нотификации

### 8.1. Назначение

HTTP-нотификации, или callback/webhook, позволяют оперативно отслеживать состояние заказа.

Когда приходит webhook, страница требует:

- проверить состояние заявки методом `Info`
- затем поменять статус заказа в системе учета интернет-магазина

Вебхуки полезны в ситуациях, когда объект API меняется без участия магазина.

### 8.2. Варианты настройки адреса webhook

#### 8.2.1. Статичный адрес

На странице описаны шаги:

1. Вписать адрес в личном кабинете руководителя.
2. Перейти в раздел «Кредитование в магазинах».
3. Выбрать магазин.
4. Нажать «Редактировать».
5. Нажать «Уведомления».
6. Заполнить «Адрес для HTTP-нотификаций».
7. Поставить галочку в чекбоксе.

#### 8.2.2. Свой адрес для каждой заявки

- Его нужно передавать при создании заявки в поле `webhookURL`.
- При этом в личном кабинете руководителя все равно должен быть настроен статичный адрес.

### 8.3. Ограничения и условия для webhook

- Домен уведомлений должен совпадать с доменом сайта магазина.
- Для использования вебхуков при создании заявки нужно передавать `orderNumber`.
- В webhook поле `id` соответствует `orderNumber`.
- Сайт должен использовать защищенный `HTTPS`.
- Решения с `HTTP` не поддерживаются.

### 8.4. События webhook

На странице перечислены события:

- `Approved` — заявка одобрена одним из банков, покупателю остается подписать документы через СМС или на встрече с представителем банка
- `Rejected` — по заявке пришел отказ от одного из банков, можно связаться с покупателем и предложить альтернативные способы оплаты
- `Canceled` — заявка отменена, покупатель по какой-то причине отменил заказ
- `Signed` — договор подписан покупателем; если заявка не охлаждается, после этого можно подтверждать заявку методом `Commit` или в личном кабинете агента, а затем выдавать товар или оказывать услугу
- `Signed` при `commit_cooldown.has_happened = true` — период охлаждения закончен

### 8.5. Дополнительные примечания по событиям webhook

- Время действия при охлаждении можно посмотреть в параметре webhook `Signed`, где `commit_cooldown.has_happened = true`.
- Если по заявке действовал период охлаждения, webhook `Signed` с `commit_cooldown.has_happened = true` будет отправлен после окончания периода охлаждения.
- При автокоммите такой webhook отправится, если заявка не была отменена клиентом.
- При отключенном автокоммите после такого webhook можно подтверждать заявку методом `Commit` или в личном кабинете агента.
- Если клиент отменит заявку в период охлаждения, будет отправлен webhook `Canceled`.
- Вебхуки с решениями `Approved` и `Rejected` приходят по каждому банку.
- Вебхук `Rejected` не означает, что по всей заявке уже пришел окончательный отказ, нужно дождаться решения всех банков.

### 8.6. Проверка подлинности webhook

Документация требует обязательно проверять подлинность HTTP-нотификаций:

- запросить состояние заявки через метод `Info API`
- проверить IP-адрес источника

Указанные на странице сетевые данные:

- маска сети: `91.194.226.00/23`
- IP первого хоста: `91.194.226.1`
- IP последнего хоста: `91.194.227.254`
- количество хостов в сети: `510`

### 8.7. Структуры webhook

#### 8.7.1. `CallbackRequestBody`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `id` | `string(64)` | required | Номер заказа. Соответствует `orderNumber` |
| `status` | `enum` | required | Тип вебхука |
| `created_at` | `string dateTime` | required | Дата и время создания заявки |
| `demo` | `boolean` | required | Флаг, является ли заявка тестовой |
| `committed` | `boolean` | required | Флаг, является ли заявка подтвержденной |
| `first_payment` | `number double` | not specified | Первоначальный взнос |
| `order_amount` | `number double` | required | Сумма заказа |
| `credit_amount` | `number double` | not specified | Сумма выдаваемого покупателю кредита |
| `transfer_amount` | `number double` | not specified | Сумма для перевода партнеру |
| `product` | `enum` | not specified | Тип продукта: `credit` или `installment_credit` |
| `term` | `integer int32` | not specified | Срок кредита в месяцах |
| `monthly_payment` | `number double` | not specified | Ежемесячный платеж |
| `first_name` | `string` | not specified | Имя покупателя |
| `last_name` | `string` | not specified | Фамилия покупателя |
| `middle_name` | `string` | not specified | Отчество покупателя |
| `phone` | `string` | not specified | Телефон покупателя |
| `loan_number` | `string` | not specified | Номер кредитного договора |
| `email` | `string` | not specified | E-mail покупателя |
| `appropriate_signing_types` | `string[]` | required | Доступные покупателю способы подписания |
| `signing_type` | `string` | not specified | Актуальный способ подписания заявки. Входит в `appropriate_signing_types` |
| `chosen_bank` | `string` | not specified | Банк, выдавший кредит. ОТП Банк, Т-Банк или МТС Банк |
| `expected_overdue_at` | `string dateTime` | required | Дата и время окончания срока действия заявки |
| `items` | `CallbackRequestBodyItem[]` | not specified | Состав заказа покупателя |
| `commit_cooldown` | `CallbackRequestBodyCommitCooldown` | not specified | Информация о периоде охлаждения |

#### 8.7.2. `CallbackRequestBodyStatusEnum`

- `signed` — заявка подписана
- `canceled` — заявка отменена покупателем
- `rejected` — заявка отклонена банком
- `approved` — заявка одобрена банком

#### 8.7.3. `CallbackRequestBodyProductEnum`

- `credit`
- `installment_credit`

#### 8.7.4. `CallbackRequestBodyAppropriateSigningTypesEnum`

- `bank`
- `sms`
- `ses`

#### 8.7.5. `CallbackRequestBodyItem`

Формальная таблица параметров webhook содержит только эти поля:

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `name` | `string(255)` | required | Название товарной позиции |
| `quantity` | `integer` | required | Количество единиц товара указанной позиции |
| `price` | `number` | required | Цена одной единицы товара |

#### 8.7.6. `CallbackRequestBodyItemExampleOnlyFields`

В JSON-примере webhook на той же странице внутри `items` дополнительно показаны поля, которых нет в формальной таблице параметров webhook:

| Поле | Тип в примере | Статус в документации | Описание |
|---|---|---|---|
| `category` | `string` | example only | Показано только в JSON-примере webhook |
| `vendorCode` | `null` или `string` | example only | Показано только в JSON-примере webhook |

#### 8.7.7. `CallbackRequestBodyCommitCooldown`

| Поле | Тип | Обязательность | Описание |
|---|---|---|---|
| `until` | `string dateTime` | required | Дата и время окончания периода охлаждения в ISO 8601 формате |
| `has_happened` | `boolean` | required | Флаг окончания охлаждения: `true` — охлаждение закончилось, `false` — заявка охлаждается |

### 8.8. Пример webhook

На странице явно указано, что клиентская информация в вебхуках приходит в маскированном виде.

```json
{
  "id": "1234567890",
  "status": "signed",
  "created_at": "2022-11-10T09:03:48.780Z",
  "demo": false,
  "committed": false,
  "first_payment": 0,
  "order_amount": 100000,
  "credit_amount": 100000,
  "transfer_amount": 100000,
  "product": "credit",
  "term": 6,
  "monthly_payment": 18155,
  "first_name": "Петр",
  "last_name": "И*****",
  "middle_name": "В*********",
  "phone": "+7912*****89",
  "loan_number": "543147",
  "email": "i****v@gmail.com",
  "appropriate_signing_types": [
    "bank",
    "sms"
  ],
  "signing_type": "sms",
  "chosen_bank": "Т-Банк",
  "expected_overdue_at": "2022-11-24T09:04:55.526449Z",
  "items": [
    {
      "name": "Телефон",
      "quantity": 1,
      "price": 80000,
      "category": "mobile",
      "vendorCode": null
    },
    {
      "name": "Наушники",
      "quantity": 1,
      "price": 20000
    }
  ]
}
```

### 8.9. Ответ на webhook

После получения webhook нужно отправить ответ.

На странице явно указано:

- ответ формата `2XX` считается успешным приемом webhook
- любой другой ответ приводит к повторной отправке webhook

### 8.10. Что про callback не описано на странице

На странице не описано:

- конкретный HTTP-метод, которым webhook отправляется
- состав HTTP-заголовков webhook
- требуется ли какое-то конкретное тело ответа от магазина
- сколько раз выполняется повторная отправка webhook
- с каким интервалом выполняется повторная отправка webhook

## 9. Полный список всех структур этого файла

Ниже перечислены все структуры без исключений, чтобы ими было проще пользоваться как оглавлением.

### 9.1. Create

- `CreateRequestBody`
- `CreateRequestBodyItem`
- `CreateRequestBodyValues`
- `CreateRequestBodyContact`
- `CreateRequestBodyFio`
- `CreateResponseBody`

### 9.2. Commit

- `CommitRequestPath`
- `CommitRequestHeaders`
- `CommitRequestBody`
- `CommitResponseBody`
- `CommitResponseBodyStatusEnum`
- `CommitResponseBodyProductEnum`
- `CommitResponseBodyAppropriateSigningTypesEnum`
- `CommitResponseBodyItem`
- `CommitResponseBodyCommitCooldown`

### 9.3. Cancel

- `CancelRequestPath`
- `CancelRequestHeaders`
- `CancelRequestBody`
- `CancelResponseBody`
- `CancelResponseBodyStatusEnum`
- `CancelResponseBodyProductEnum`
- `CancelResponseBodyAppropriateSigningTypesEnum`
- `CancelResponseBodyItem`
- `CancelResponseBodyCommitCooldown`

### 9.4. Info

- `InfoRequestPath`
- `InfoRequestHeaders`
- `InfoRequestBody`
- `InfoResponseBody`
- `InfoResponseBodyStatusEnum`
- `InfoResponseBodyProductEnum`
- `InfoResponseBodyAppropriateSigningTypesEnum`
- `InfoResponseBodyItem`
- `InfoResponseBodyCommitCooldown`

### 9.5. Callback

- `CallbackRequestBody`
- `CallbackRequestBodyStatusEnum`
- `CallbackRequestBodyProductEnum`
- `CallbackRequestBodyAppropriateSigningTypesEnum`
- `CallbackRequestBodyItem`
- `CallbackRequestBodyItemExampleOnlyFields`
- `CallbackRequestBodyCommitCooldown`

## 10. Полный список зафиксированных пробелов исходной документации

Этот раздел нужен специально для прозрачности: ниже перечислено все, что исходные страницы не формализуют полностью.

- Нет структуры тела ошибки ни для `Create`, ни для `Commit`, ни для `Cancel`, ни для `Info`.
- Для `Commit`, `Cancel`, `Info` и `Callback` у части полей ответа обязательность на странице не проставлена явно; в этом файле такие поля помечены как `not specified`.
- Для `Commit`, `Cancel` и `Info` тело запроса на странице не описано.
- В примере `Create` есть поле `demoFlow`, но в таблице параметров `Create` оно не задокументировано.
- В формальной таблице callback `items` содержат только `name`, `quantity`, `price`, но в JSON-примере callback дополнительно есть `category` и `vendorCode`.
- Для callback не описаны HTTP-метод, заголовки, точное тело ответа и политика retry.
