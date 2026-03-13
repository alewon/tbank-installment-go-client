package tbankinstallment

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "inprogress"
	OrderStatusApproved   OrderStatus = "approved"
	OrderStatusSigned     OrderStatus = "signed"
	OrderStatusCanceled   OrderStatus = "canceled"
	OrderStatusRejected   OrderStatus = "rejected"
)

type ProductType string

const (
	ProductTypeCredit            ProductType = "credit"
	ProductTypeInstallmentCredit ProductType = "installment_credit"
)

type SigningType string

const (
	SigningTypeBank SigningType = "bank"
	SigningTypeSMS  SigningType = "sms"
	SigningTypeSES  SigningType = "ses"
)

type CreateRequestBody struct {
	ShopID      string                  `json:"shopId"`
	ShowcaseID  string                  `json:"showcaseId"`
	Sum         float64                 `json:"sum"`
	Items       []CreateRequestBodyItem `json:"items"`
	OrderNumber string                  `json:"orderNumber,omitempty"`
	PromoCode   string                  `json:"promoCode,omitempty"`
	WebhookURL  string                  `json:"webhookURL,omitempty"`
	SuccessURL  string                  `json:"successURL,omitempty"`
	FailURL     string                  `json:"failURL,omitempty"`
	ReturnURL   string                  `json:"returnURL,omitempty"`
	Values      *CreateRequestValues    `json:"values,omitempty"`
}

type CreateRequestBodyItem struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Category   string  `json:"category,omitempty"`
	VendorCode string  `json:"vendorCode,omitempty"`
}

type CreateRequestValues struct {
	Contact *CreateContact `json:"contact,omitempty"`
}

type CreateRequestBodyValues = CreateRequestValues

type CreateContact struct {
	FIO         *CreateFIO `json:"fio,omitempty"`
	MobilePhone string     `json:"mobilePhone,omitempty"`
	Email       string     `json:"email,omitempty"`
}

type CreateRequestBodyContact = CreateContact

type CreateFIO struct {
	LastName   string `json:"lastName,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
}

type CreateRequestBodyFio = CreateFIO

type CreateResponse struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type CreateRequestBodyExampleOnlyFields struct {
	DemoFlow string `json:"demoFlow,omitempty"`
}

type CreateRequest = CreateRequestBody

type CreateItem = CreateRequestBodyItem

type CreateResponseBody = CreateResponse

type CommitRequestPath struct {
	OrderNumber string `json:"orderNumber"`
}

type CommitRequestHeaders struct {
	Authorization string `json:"Authorization"`
}

type CommitRequestBody struct{}

type CommitResponseBody struct {
	ID                      string                            `json:"id"`
	Status                  OrderStatus                       `json:"status"`
	CreatedAt               time.Time                         `json:"created_at"`
	Demo                    bool                              `json:"demo"`
	Committed               bool                              `json:"committed"`
	FirstPayment            *float64                          `json:"first_payment,omitempty"`
	OrderAmount             float64                           `json:"order_amount"`
	CreditAmount            *float64                          `json:"credit_amount,omitempty"`
	TransferAmount          *float64                          `json:"transfer_amount,omitempty"`
	Product                 ProductType                       `json:"product,omitempty"`
	Term                    *int                              `json:"term,omitempty"`
	MonthlyPayment          *float64                          `json:"monthly_payment,omitempty"`
	FirstName               string                            `json:"first_name,omitempty"`
	LastName                string                            `json:"last_name,omitempty"`
	MiddleName              string                            `json:"middle_name,omitempty"`
	Phone                   string                            `json:"phone,omitempty"`
	LoanNumber              string                            `json:"loan_number,omitempty"`
	Email                   string                            `json:"email,omitempty"`
	AppropriateSigningTypes []SigningType                     `json:"appropriate_signing_types"`
	SigningType             SigningType                       `json:"signing_type,omitempty"`
	ChosenBank              string                            `json:"chosen_bank,omitempty"`
	ExpectedOverdueAt       time.Time                         `json:"expected_overdue_at"`
	Items                   []CommitResponseBodyItem          `json:"items,omitempty"`
	CommitCooldown          *CommitResponseBodyCommitCooldown `json:"commit_cooldown,omitempty"`
}

type CommitResponseBodyItem struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Category   string  `json:"category,omitempty"`
	VendorCode string  `json:"vendorCode,omitempty"`
}

type CommitResponseBodyCommitCooldown struct {
	Until time.Time `json:"until"`
}

type CommitResponseBodyStatusEnum = OrderStatus

type CommitResponseBodyProductEnum = ProductType

type CommitResponseBodyAppropriateSigningTypesEnum = SigningType

type CommitResponse = CommitResponseBody

type CommitResponseItem = CommitResponseBodyItem

type CommitResponseCooldown = CommitResponseBodyCommitCooldown

type CancelRequestPath struct {
	OrderNumber string `json:"orderNumber"`
}

type CancelRequestHeaders struct {
	Authorization string `json:"Authorization"`
}

type CancelRequestBody struct{}

type CancelResponseBody struct {
	ID                      string                            `json:"id"`
	Status                  OrderStatus                       `json:"status"`
	CreatedAt               time.Time                         `json:"created_at"`
	Demo                    bool                              `json:"demo"`
	Committed               bool                              `json:"committed"`
	FirstPayment            *float64                          `json:"first_payment,omitempty"`
	OrderAmount             float64                           `json:"order_amount"`
	CreditAmount            *float64                          `json:"credit_amount,omitempty"`
	TransferAmount          *float64                          `json:"transfer_amount,omitempty"`
	Product                 ProductType                       `json:"product,omitempty"`
	Term                    *int                              `json:"term,omitempty"`
	MonthlyPayment          *float64                          `json:"monthly_payment,omitempty"`
	FirstName               string                            `json:"first_name,omitempty"`
	LastName                string                            `json:"last_name,omitempty"`
	MiddleName              string                            `json:"middle_name,omitempty"`
	Phone                   string                            `json:"phone,omitempty"`
	LoanNumber              string                            `json:"loan_number,omitempty"`
	Email                   string                            `json:"email,omitempty"`
	AppropriateSigningTypes []SigningType                     `json:"appropriate_signing_types"`
	SigningType             SigningType                       `json:"signing_type,omitempty"`
	ChosenBank              string                            `json:"chosen_bank,omitempty"`
	ExpectedOverdueAt       time.Time                         `json:"expected_overdue_at"`
	Items                   []CancelResponseBodyItem          `json:"items,omitempty"`
	CommitCooldown          *CancelResponseBodyCommitCooldown `json:"commit_cooldown,omitempty"`
}

type CancelResponseBodyItem struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Category   string  `json:"category,omitempty"`
	VendorCode string  `json:"vendorCode,omitempty"`
}

type CancelResponseBodyCommitCooldown struct {
	Until time.Time `json:"until"`
}

type CancelResponseBodyStatusEnum = OrderStatus

type CancelResponseBodyProductEnum = ProductType

type CancelResponseBodyAppropriateSigningTypesEnum = SigningType

type CancelResponse = CancelResponseBody

type CancelResponseItem = CancelResponseBodyItem

type CancelResponseCooldown = CancelResponseBodyCommitCooldown

type InfoRequestPath struct {
	OrderNumber string `json:"orderNumber"`
}

type InfoRequestHeaders struct {
	Authorization string `json:"Authorization"`
}

type InfoRequestBody struct{}

type InfoResponseBody struct {
	ID                      string                          `json:"id"`
	Status                  OrderStatus                     `json:"status"`
	CreatedAt               time.Time                       `json:"created_at"`
	Demo                    bool                            `json:"demo"`
	Committed               bool                            `json:"committed"`
	FirstPayment            *float64                        `json:"first_payment,omitempty"`
	OrderAmount             float64                         `json:"order_amount"`
	CreditAmount            *float64                        `json:"credit_amount,omitempty"`
	TransferAmount          *float64                        `json:"transfer_amount,omitempty"`
	Product                 ProductType                     `json:"product,omitempty"`
	Term                    *int                            `json:"term,omitempty"`
	MonthlyPayment          *float64                        `json:"monthly_payment,omitempty"`
	FirstName               string                          `json:"first_name,omitempty"`
	LastName                string                          `json:"last_name,omitempty"`
	MiddleName              string                          `json:"middle_name,omitempty"`
	Phone                   string                          `json:"phone,omitempty"`
	LoanNumber              string                          `json:"loan_number,omitempty"`
	Email                   string                          `json:"email,omitempty"`
	AppropriateSigningTypes []SigningType                   `json:"appropriate_signing_types"`
	SigningType             SigningType                     `json:"signing_type,omitempty"`
	ChosenBank              string                          `json:"chosen_bank,omitempty"`
	ExpectedOverdueAt       time.Time                       `json:"expected_overdue_at"`
	Items                   []InfoResponseBodyItem          `json:"items,omitempty"`
	CommitCooldown          *InfoResponseBodyCommitCooldown `json:"commit_cooldown,omitempty"`
}

type InfoResponseBodyItem struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Category   string  `json:"category,omitempty"`
	VendorCode string  `json:"vendorCode,omitempty"`
}

type InfoResponseBodyCommitCooldown struct {
	Until time.Time `json:"until"`
}

type InfoResponseBodyStatusEnum = OrderStatus

type InfoResponseBodyProductEnum = ProductType

type InfoResponseBodyAppropriateSigningTypesEnum = SigningType

type InfoResponse = InfoResponseBody

type InfoResponseItem = InfoResponseBodyItem

type InfoResponseCooldown = InfoResponseBodyCommitCooldown

type CallbackRequestBodyStatusEnum string

const (
	CallbackStatusSigned   CallbackRequestBodyStatusEnum = "signed"
	CallbackStatusCanceled CallbackRequestBodyStatusEnum = "canceled"
	CallbackStatusRejected CallbackRequestBodyStatusEnum = "rejected"
	CallbackStatusApproved CallbackRequestBodyStatusEnum = "approved"
)

type CallbackRequestBody struct {
	ID                      string                             `json:"id"`
	Status                  CallbackRequestBodyStatusEnum      `json:"status"`
	CreatedAt               time.Time                          `json:"created_at"`
	Demo                    bool                               `json:"demo"`
	Committed               bool                               `json:"committed"`
	FirstPayment            *float64                           `json:"first_payment,omitempty"`
	OrderAmount             float64                            `json:"order_amount"`
	CreditAmount            *float64                           `json:"credit_amount,omitempty"`
	TransferAmount          *float64                           `json:"transfer_amount,omitempty"`
	Product                 ProductType                        `json:"product,omitempty"`
	Term                    *int                               `json:"term,omitempty"`
	MonthlyPayment          *float64                           `json:"monthly_payment,omitempty"`
	FirstName               string                             `json:"first_name,omitempty"`
	LastName                string                             `json:"last_name,omitempty"`
	MiddleName              string                             `json:"middle_name,omitempty"`
	Phone                   string                             `json:"phone,omitempty"`
	LoanNumber              string                             `json:"loan_number,omitempty"`
	Email                   string                             `json:"email,omitempty"`
	AppropriateSigningTypes []SigningType                      `json:"appropriate_signing_types"`
	SigningType             SigningType                        `json:"signing_type,omitempty"`
	ChosenBank              string                             `json:"chosen_bank,omitempty"`
	ExpectedOverdueAt       time.Time                          `json:"expected_overdue_at"`
	Items                   []CallbackRequestBodyItem          `json:"items,omitempty"`
	CommitCooldown          *CallbackRequestBodyCommitCooldown `json:"commit_cooldown,omitempty"`
}

type CallbackRequestBodyItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	CallbackRequestBodyItemExampleOnlyFields
}

type CallbackRequestBodyItemExampleOnlyFields struct {
	Category   string  `json:"category,omitempty"`
	VendorCode *string `json:"vendorCode,omitempty"`
}

type CallbackRequestBodyCommitCooldown struct {
	Until       time.Time `json:"until"`
	HasHappened bool      `json:"has_happened"`
}

type CallbackRequestBodyProductEnum = ProductType

type CallbackRequestBodyAppropriateSigningTypesEnum = SigningType

type CallbackStatus = CallbackRequestBodyStatusEnum

type WebhookPayload = CallbackRequestBody

type WebhookItem = CallbackRequestBodyItem

type WebhookCooldown = CallbackRequestBodyCommitCooldown

func ParseWebhook(r io.Reader) (*CallbackRequestBody, error) {
	var payload CallbackRequestBody
	if err := json.NewDecoder(r).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode webhook payload: %w", err)
	}
	return &payload, nil
}
