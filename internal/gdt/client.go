// Package gdt implements a client for the GDT (General Department of
// Taxation) e-invoice API. Wire contract (isolated here — GDT's real
// SOAP/HTTPS endpoints may drift; swap formats without touching callers):
//
//	POST {base}/api/invoice/submit   {"xml": "...", "cert_id": "..."}
//	      200 → {"transaction_id","status","gdt_ref","message"}
//	GET  {base}/api/invoice/status?transaction_id=...
//	      200 → {"status","gdt_ref","message"}
//	POST {base}/api/invoice/cancel   {"transaction_id": "...", "reason": "..."}
//	      200 → 200
//	POST {base}/api/submission/declare   {"xml": "...", "cert_id": "..."}
//	      200 → {"submission_id","status","ack_ref","message"}
//	GET  {base}/api/submission/status?id=...
//	      200 → {"status","ack_ref","message"}
//
// Retry policy: 5xx and network errors are retried with backoff
// (1s/5s/30s by default, max 3 attempts). 4xx are terminal.
package gdt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	// ErrUpstream — GDT unreachable or persistently failing (5xx/network).
	ErrUpstream = errors.New("gdt: upstream error")
	// ErrUnauthorized — 401, no retry.
	ErrUnauthorized = errors.New("gdt: unauthorized")
	// ErrInvalidRequest — 400/422, no retry.
	ErrInvalidRequest = errors.New("gdt: invalid request")
	// ErrNotFound — 404, no retry.
	ErrNotFound = errors.New("gdt: not found")
)

// Client is a GDT API client.
type Client struct {
	baseURL   string
	http      *http.Client
	retries   []time.Duration
	authToken string
}

// Option configures a Client.
type Option func(*Client)

// WithTimeout sets the per-request timeout (default 120s).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.http.Timeout = d }
}

// WithRetry sets retry backoff delays after the first attempt (default 1s/5s/30s).
func WithRetry(delays ...time.Duration) Option {
	return func(c *Client) { c.retries = delays }
}

// WithToken sets a static bearer token on every request.
func WithToken(token string) Option {
	return func(c *Client) { c.authToken = token }
}

// New builds a Client. url must be an absolute http(s) base URL.
func New(baseURL string, opts ...Option) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("gdt: invalid base URL %q", baseURL)
	}
	c := &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 120 * time.Second},
		retries: []time.Duration{time.Second, 5 * time.Second, 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// SubmitRequest is the invoice submission payload.
type SubmitRequest struct {
	XML    string `json:"xml"`
	CertID string `json:"cert_id"`
}

// SubmitResponse is the GDT acknowledgement of a submission.
type SubmitResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	GDTRef        string `json:"gdt_ref"`
	Message       string `json:"message,omitempty"`
}

// StatusResponse is the invoice status payload.
type StatusResponse struct {
	Status  string `json:"status"`
	GDTRef  string `json:"gdt_ref,omitempty"`
	Message string `json:"message,omitempty"`
}

// SubmitInvoice submits a signed invoice XML to GDT.
func (c *Client) SubmitInvoice(ctx context.Context, invoiceXML, certID string) (*SubmitResponse, error) {
	body, _ := json.Marshal(SubmitRequest{XML: invoiceXML, CertID: certID})
	var out SubmitResponse
	if err := c.do(ctx, http.MethodPost, "/api/invoice/submit", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetInvoiceStatus polls GDT for invoice status.
func (c *Client) GetInvoiceStatus(ctx context.Context, transactionID string) (*StatusResponse, error) {
	var out StatusResponse
	if err := c.do(ctx, http.MethodGet,
		"/api/invoice/status?transaction_id="+url.QueryEscape(transactionID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelInvoice requests cancellation of a submitted invoice.
func (c *Client) CancelInvoice(ctx context.Context, transactionID, reason string) error {
	body, _ := json.Marshal(map[string]string{"transaction_id": transactionID, "reason": reason})
	return c.do(ctx, http.MethodPost, "/api/invoice/cancel", body, nil)
}

// DeclarationSubmitResponse is the GDT acknowledgement of a declaration
// submission.
type DeclarationSubmitResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
	// Code is the GDT response code (spec §4.2): 00 acknowledged, 01 schema
	// failure, 02 duplicate, 03 tax code not found, 10 period already
	// declared, 99 system error.
	Code    string `json:"code,omitempty"`
	AckRef  string `json:"ack_ref,omitempty"`
	Message string `json:"message,omitempty"`
}

// DeclarationStatusResponse is the declaration status payload.
type DeclarationStatusResponse struct {
	Status string `json:"status"`
	// Code carries the same GDT response codes as DeclarationSubmitResponse.
	Code    string `json:"code,omitempty"`
	AckRef  string `json:"ack_ref,omitempty"`
	Message string `json:"message,omitempty"`
}

// SubmitDeclaration submits a signed declaration XML (HTKK file) to GDT.
func (c *Client) SubmitDeclaration(ctx context.Context, declarationXML, certID string) (*DeclarationSubmitResponse, error) {
	body, _ := json.Marshal(SubmitRequest{XML: declarationXML, CertID: certID})
	var out DeclarationSubmitResponse
	if err := c.do(ctx, http.MethodPost, "/api/submission/declare", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// QueryDeclarationStatus polls GDT for declaration status.
func (c *Client) QueryDeclarationStatus(ctx context.Context, submissionID string) (*DeclarationStatusResponse, error) {
	var out DeclarationStatusResponse
	if err := c.do(ctx, http.MethodGet,
		"/api/submission/status?id="+url.QueryEscape(submissionID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// do performs a request with retry on 5xx/network errors.
func (c *Client) do(ctx context.Context, method, path string, body []byte, out any) error {
	var lastErr error
	attempts := 1 + len(c.retries)
	for i := 0; i < attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		var reader io.Reader
		if body != nil {
			reader = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
		if err != nil {
			return fmt.Errorf("gdt: build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if c.authToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.authToken)
		}
		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = ErrUpstream // network error → retry
		} else {
			respBody, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			switch {
			case resp.StatusCode >= 200 && resp.StatusCode < 300:
				if readErr != nil {
					return fmt.Errorf("gdt: read response: %w", readErr)
				}
				if out != nil && len(respBody) > 0 {
					if err := json.Unmarshal(respBody, out); err != nil {
						return fmt.Errorf("gdt: decode response: %w", err)
					}
				}
				return nil
			case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
				return ErrUnauthorized
			case resp.StatusCode == http.StatusNotFound:
				return ErrNotFound
			case resp.StatusCode >= 400 && resp.StatusCode < 500:
				return ErrInvalidRequest
			default: // 5xx → retry
				lastErr = ErrUpstream
			}
		}
		if i < len(c.retries) {
			select {
			case <-time.After(c.retries[i]):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return lastErr
}
