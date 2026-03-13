package tbankinstallment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://forma.tbank.ru"
const defaultPromoCode = "default"

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	Username             string
	Password             string
	BaseURL              string
	HTTPClient           HTTPDoer
	UserAgent            string
	WebhookTrustedSubnet string
}

type Client struct {
	baseURL              *url.URL
	username             string
	password             string
	httpClient           HTTPDoer
	userAgent            string
	webhookTrustedSubnet *net.IPNet
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.Username) == "" {
		return nil, fmt.Errorf("username is required")
	}
	if strings.TrimSpace(cfg.Password) == "" {
		return nil, fmt.Errorf("password is required")
	}

	rawBaseURL := cfg.BaseURL
	if rawBaseURL == "" {
		rawBaseURL = defaultBaseURL
	}

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("base url must include scheme and host")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	var webhookTrustedSubnet *net.IPNet
	if strings.TrimSpace(cfg.WebhookTrustedSubnet) != "" {
		_, network, err := net.ParseCIDR(cfg.WebhookTrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("parse webhook trusted subnet: %w", err)
		}
		webhookTrustedSubnet = network
	}

	return &Client{
		baseURL:              baseURL,
		username:             cfg.Username,
		password:             cfg.Password,
		httpClient:           httpClient,
		userAgent:            cfg.UserAgent,
		webhookTrustedSubnet: webhookTrustedSubnet,
	}, nil
}

func (c *Client) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if strings.TrimSpace(req.PromoCode) == "" {
		req.PromoCode = defaultPromoCode
	}

	var resp CreateResponse
	if err := c.doJSON(ctx, http.MethodPost, "/api/partners/v2/orders/create", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) Commit(ctx context.Context, orderNumber string) (*CommitResponse, error) {
	if strings.TrimSpace(orderNumber) == "" {
		return nil, fmt.Errorf("orderNumber is required")
	}

	var resp CommitResponse
	if err := c.doJSON(ctx, http.MethodPost, orderActionPath(orderNumber, "commit"), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) Cancel(ctx context.Context, orderNumber string) (*CancelResponse, error) {
	if strings.TrimSpace(orderNumber) == "" {
		return nil, fmt.Errorf("orderNumber is required")
	}

	var resp CancelResponse
	if err := c.doJSON(ctx, http.MethodPost, orderActionPath(orderNumber, "cancel"), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) Info(ctx context.Context, orderNumber string) (*InfoResponse, error) {
	if strings.TrimSpace(orderNumber) == "" {
		return nil, fmt.Errorf("orderNumber is required")
	}

	var resp InfoResponse
	if err := c.doJSON(ctx, http.MethodGet, orderActionPath(orderNumber, "info"), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func orderActionPath(orderNumber, action string) string {
	return "/api/partners/v2/orders/" + url.PathEscape(strings.TrimSpace(orderNumber)) + "/" + action
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, dst any) error {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, req.URL.String(), err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp, payload)
	}
	if dst == nil || len(payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(payload, dst); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse request path: %w", err)
	}

	fullURL := c.baseURL.ResolveReference(rel)

	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
		reader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

type APIError struct {
	StatusCode int
	Status     string
	RetryAfter time.Duration
	Body       []byte
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.RetryAfter > 0 {
		return fmt.Sprintf("tbank api error: %s (retry after %s)", e.Status, e.RetryAfter)
	}
	return fmt.Sprintf("tbank api error: %s", e.Status)
}

func newAPIError(resp *http.Response, body []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       append([]byte(nil), body...),
	}

	if retryAfter := strings.TrimSpace(resp.Header.Get("Retry-After")); retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds >= 0 {
			apiErr.RetryAfter = time.Duration(seconds) * time.Second
		}
	}

	return apiErr
}

func (c *Client) IsTrustedWebhookIP(ip net.IP) bool {
	return c != nil && c.webhookTrustedSubnet != nil && ip != nil && c.webhookTrustedSubnet.Contains(ip)
}

func (c *Client) IsTrustedWebhookRequest(r *http.Request) bool {
	if c == nil || r == nil {
		return false
	}

	ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ipStr = r.RemoteAddr
	}

	return c.IsTrustedWebhookIP(net.ParseIP(ipStr))
}
