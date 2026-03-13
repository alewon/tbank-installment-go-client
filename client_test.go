package tbankinstallment

import (
	"context"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestCreateDoesNotUseBasicAuthAndBuildsPayload(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/api/partners/v2/orders/create" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}

			if got := req.Header.Get("Authorization"); got != "" {
				t.Fatalf("create request must not include auth header: %s", got)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !strings.Contains(string(body), `"shopId":"shop"`) {
				t.Fatalf("unexpected body: %s", string(body))
			}
			if !strings.Contains(string(body), `"promoCode":"default"`) {
				t.Fatalf("unexpected body: %s", string(body))
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":"1","link":"https://forma.example/app/1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Create(context.Background(), CreateRequest{
		ShopID:     "shop",
		ShowcaseID: "showcase",
		Sum:        1000,
		Items: []CreateItem{{
			Name:     "item",
			Quantity: 1,
			Price:    1000,
		}},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if resp.ID != "1" {
		t.Fatalf("unexpected response id: %s", resp.ID)
	}
}

func TestCreateDemoUsesCreateDemoEndpoint(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/api/partners/v2/orders/create-demo" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}
			if got := req.Header.Get("Authorization"); got != "" {
				t.Fatalf("create demo request must not include auth header: %s", got)
			}

			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			raw := string(body)
			for _, fragment := range []string{
				`"promoCode":"default"`,
				`"demoFlow":"sms"`,
			} {
				if !strings.Contains(raw, fragment) {
					t.Fatalf("request body does not contain %s: %s", fragment, raw)
				}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":"1","link":"https://forma.example/app/1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.CreateDemo(context.Background(), CreateDemoRequest{
		CreateRequestBody: CreateRequestBody{
			ShopID:     "shop",
			ShowcaseID: "showcase",
			Sum:        1000,
			Items: []CreateItem{{
				Name:     "item",
				Quantity: 1,
				Price:    1000,
			}},
		},
		DemoFlow: DemoFlowSMS,
	})
	if err != nil {
		t.Fatalf("CreateDemo() error = %v", err)
	}
	if resp.ID != "1" {
		t.Fatalf("unexpected response id: %s", resp.ID)
	}
}

func TestCreateEncodesPrefilledCustomerValues(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			raw := string(body)
			for _, fragment := range []string{
				`"values":{"contact":{"fio":{"lastName":"Иванов","firstName":"Иван","middleName":"Иванович"},"mobilePhone":"9998887766","email":"ivan@example.com"}}`,
				`"orderNumber":"1234567890"`,
				`"webhookURL":"https://shop.example/webhook"`,
				`"successURL":"https://shop.example/success"`,
				`"failURL":"https://shop.example/fail"`,
				`"returnURL":"https://shop.example/return"`,
			} {
				if !strings.Contains(raw, fragment) {
					t.Fatalf("request body does not contain %s: %s", fragment, raw)
				}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":"1","link":"https://forma.example/app/1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.Create(context.Background(), CreateRequest{
		ShopID:      "shop",
		ShowcaseID:  "showcase",
		Sum:         1000,
		OrderNumber: "1234567890",
		WebhookURL:  "https://shop.example/webhook",
		SuccessURL:  "https://shop.example/success",
		FailURL:     "https://shop.example/fail",
		ReturnURL:   "https://shop.example/return",
		Items: []CreateItem{{
			Name:     "item",
			Quantity: 1,
			Price:    1000,
		}},
		Values: &CreateRequestValues{
			Contact: &CreateContact{
				FIO: &CreateFIO{
					LastName:   "Иванов",
					FirstName:  "Иван",
					MiddleName: "Иванович",
				},
				MobilePhone: "9998887766",
				Email:       "ivan@example.com",
			},
		},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestCreateOmitsEmptyReturnURLs(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			raw := string(body)
			for _, fragment := range []string{
				`"successURL"`,
				`"failURL"`,
				`"returnURL"`,
			} {
				if strings.Contains(raw, fragment) {
					t.Fatalf("request body must not contain %s: %s", fragment, raw)
				}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":"1","link":"https://forma.example/app/1"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.Create(context.Background(), CreateRequest{
		ShopID:     "shop",
		ShowcaseID: "showcase",
		Sum:        1000,
		Items: []CreateItem{{
			Name:     "item",
			Quantity: 1,
			Price:    1000,
		}},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestCommitUsesPOSTAndDecodesCommitResponse(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/api/partners/v2/orders/order-42/commit" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}
			wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))
			if got := req.Header.Get("Authorization"); got != wantAuth {
				t.Fatalf("unexpected auth header: %s", got)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body: io.NopCloser(strings.NewReader(`{
					"id":"order-42",
					"status":"signed",
					"created_at":"2022-11-10T09:03:48.780Z",
					"demo":false,
					"committed":true,
					"order_amount":100000,
					"appropriate_signing_types":["sms","bank"],
					"expected_overdue_at":"2022-11-24T09:04:55.526449Z",
					"commit_cooldown":{"until":"2022-11-24T09:04:55.526449Z"}
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Commit(context.Background(), "order-42")
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	if resp.ID != "order-42" || !resp.Committed || resp.Status != OrderStatusSigned {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.CommitCooldown == nil {
		t.Fatal("expected commit cooldown to be decoded")
	}
}

func TestInfoUsesDemoPrefixedBasicAuthInDemoMode(t *testing.T) {
	client, err := NewClient(Config{
		Username: "showcase",
		Password: "pass",
		Demo:     true,
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("demo-showcase:pass"))
			if got := req.Header.Get("Authorization"); got != wantAuth {
				t.Fatalf("unexpected auth header: %s", got)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body: io.NopCloser(strings.NewReader(`{
					"id":"order-42",
					"status":"signed",
					"created_at":"2022-11-10T09:03:48.780Z",
					"demo":true,
					"committed":false,
					"order_amount":3200,
					"appropriate_signing_types":["sms"],
					"expected_overdue_at":"2022-11-24T09:04:55.526449Z"
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Info(context.Background(), "order-42")
	if err != nil {
		t.Fatalf("Info() error = %v", err)
	}
	if resp.ID != "order-42" || !resp.Demo {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCancelUsesPOSTAndDecodesCancelResponse(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", req.Method)
			}
			if req.URL.Path != "/api/partners/v2/orders/order-42/cancel" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body: io.NopCloser(strings.NewReader(`{
					"id":"order-42",
					"status":"canceled",
					"created_at":"2022-11-10T09:03:48.780Z",
					"demo":false,
					"committed":false,
					"order_amount":100000,
					"appropriate_signing_types":["sms"],
					"expected_overdue_at":"2022-11-24T09:04:55.526449Z"
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Cancel(context.Background(), "order-42")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if resp.ID != "order-42" || resp.Status != OrderStatusCanceled {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestInfoEscapesOrderNumber(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.EscapedPath() != "/api/partners/v2/orders/order%2F42/info" {
				t.Fatalf("unexpected escaped path: %s", req.URL.EscapedPath())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body: io.NopCloser(strings.NewReader(`{
					"id":"order/42",
					"status":"approved",
					"created_at":"2022-11-10T09:03:48.780Z",
					"demo":false,
					"committed":false,
					"order_amount":100000,
					"appropriate_signing_types":["sms"],
					"expected_overdue_at":"2022-11-24T09:04:55.526449Z"
				}`)),
				Header: make(http.Header),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Info(context.Background(), "order/42")
	if err != nil {
		t.Fatalf("Info() error = %v", err)
	}
	if resp.Status != OrderStatusApproved {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
}

func TestAPIErrorParsesRetryAfter(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			header := make(http.Header)
			header.Set("Retry-After", "3")
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Status:     "429 Too Many Requests",
				Body:       io.NopCloser(strings.NewReader(`rate limit`)),
				Header:     header,
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.Info(context.Background(), "order")
	if err == nil {
		t.Fatal("Info() error = nil, want error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
	if apiErr.RetryAfter != 3*time.Second {
		t.Fatalf("unexpected retry after: %s", apiErr.RetryAfter)
	}
}

func TestPathMethodsRequireOrderNumber(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.Commit(context.Background(), ""); err == nil {
		t.Fatal("Commit() error = nil, want error")
	}
	if _, err := client.Cancel(context.Background(), ""); err == nil {
		t.Fatal("Cancel() error = nil, want error")
	}
	if _, err := client.Info(context.Background(), ""); err == nil {
		t.Fatal("Info() error = nil, want error")
	}
}

func TestParseWebhook(t *testing.T) {
	payload, err := ParseWebhook(strings.NewReader(`{
			"id":"1234567890",
			"status":"signed",
		"created_at":"2022-11-10T09:03:48.780Z",
		"demo":false,
		"committed":false,
		"order_amount":100000,
		"appropriate_signing_types":["bank","sms"],
		"expected_overdue_at":"2022-11-24T09:04:55.526449Z",
		"commit_cooldown":{"until":"2022-11-24T09:04:55.526449Z","has_happened":true},
			"items":[{"name":"Телефон","quantity":1,"price":80000,"category":"mobile","vendorCode":null}]
		}`))
	if err != nil {
		t.Fatalf("ParseWebhook() error = %v", err)
	}
	if payload.Status != CallbackStatusSigned {
		t.Fatalf("unexpected status: %s", payload.Status)
	}
	if payload.CommitCooldown == nil || !payload.CommitCooldown.HasHappened {
		t.Fatal("expected commit cooldown with has_happened = true")
	}
	if len(payload.Items) != 1 {
		t.Fatalf("unexpected items count: %d", len(payload.Items))
	}
	if payload.Items[0].Category != "mobile" {
		t.Fatalf("unexpected item category: %s", payload.Items[0].Category)
	}
	if payload.Items[0].VendorCode != nil {
		t.Fatalf("expected nil vendorCode, got %v", *payload.Items[0].VendorCode)
	}
}

func TestIsTrustedWebhookRequest(t *testing.T) {
	client, err := NewClient(Config{
		Username:             "user",
		Password:             "pass",
		BaseURL:              "https://example.com",
		WebhookTrustedSubnet: "91.194.226.0/23",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := &http.Request{RemoteAddr: "91.194.226.10:12345"}
	if !client.IsTrustedWebhookRequest(req) {
		t.Fatal("expected request to be trusted")
	}
	if client.IsTrustedWebhookIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("unexpected trusted ip")
	}
}

func TestTrustedWebhookRequestWithoutConfiguredSubnet(t *testing.T) {
	client, err := NewClient(Config{
		Username: "user",
		Password: "pass",
		BaseURL:  "https://example.com",
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.IsTrustedWebhookRequest(&http.Request{RemoteAddr: "91.194.226.10:12345"}) {
		t.Fatal("expected request to be untrusted when subnet is not configured")
	}
}

func TestNewClientRejectsInvalidWebhookTrustedSubnet(t *testing.T) {
	_, err := NewClient(Config{
		Username:             "user",
		Password:             "pass",
		BaseURL:              "https://example.com",
		WebhookTrustedSubnet: "not-a-cidr",
	})
	if err == nil {
		t.Fatal("NewClient() error = nil, want error")
	}
}
