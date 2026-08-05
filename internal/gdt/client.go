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
//	      200 → {"submission_id","status","code","ack_ref","message"}
//	GET  {base}/api/submission/status?id=...
//	      200 → {"status","code","ack_ref","message"}
//
// Errors are mapped to domain errors at this boundary so callers never see
// transport details: 401/403 → ErrGDTUnauthorized, 404/4xx → ErrGDTRejected,
// 5xx/network → ErrGDTUnavailable.
//
// Retry policy: 5xx and network errors are retried with backoff
// (1s/5s/30s by default, max 3 attempts). 4xx are terminal.
package gdt

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"gotax/internal/domain"
)

// GDTError is a structured error from the GDT API. It wraps domain errors
// with HTTP-level details for debugging and audit trails.
type GDTError struct {
	Err        error  // underlying domain error (ErrGDTRejected, etc.)
	StatusCode int    // HTTP status code
	Body       string // raw response body (truncated to 4KB)
	Path       string // request path
}

func (e *GDTError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("gdt: %s (HTTP %d): %s", e.Path, e.StatusCode, e.Body)
	}
	return fmt.Sprintf("gdt: %s (HTTP %d)", e.Path, e.StatusCode)
}

func (e *GDTError) Unwrap() error { return e.Err }

// Client is a GDT API client.
type Client struct {
	baseURL   string
	http      *http.Client
	retries   []time.Duration
	authToken string
	logger    func(string, ...any)
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

// WithLogger sets a logger for request/response debugging.
func WithLogger(l func(string, ...any)) Option {
	return func(c *Client) { c.logger = l }
}

// WithClientCert configures mTLS with a client certificate and CA bundle.
// certPath: PEM file with client cert + key. caPath: PEM file with CA certs (empty = system roots).
func WithClientCert(certPath, caPath string) Option {
	return func(c *Client) {
		cert, err := tls.LoadX509KeyPair(certPath, certPath)
		if err != nil {
			if c.logger != nil {
				c.logger("gdt: load client cert: %v", err)
			}
			return
		}
		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		if caPath != "" {
			ca, err := os.ReadFile(caPath)
			if err != nil {
				if c.logger != nil {
					c.logger("gdt: read CA bundle: %v", err)
				}
				return
			}
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM(ca) {
				if c.logger != nil {
					c.logger("gdt: parse CA bundle: no valid certs")
				}
				return
			}
			tlsCfg.RootCAs = pool
		}
		c.http.Transport = &http.Transport{TLSClientConfig: tlsCfg}
	}
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

// SubmitInvoice submits a signed invoice XML to GDT.
func (c *Client) SubmitInvoice(ctx context.Context, invoiceXML, certID string) (*domain.GDTSubmitResponse, error) {
	body, _ := json.Marshal(SubmitRequest{XML: invoiceXML, CertID: certID})
	var out domain.GDTSubmitResponse
	if err := c.do(ctx, http.MethodPost, "/api/invoice/submit", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetInvoiceStatus polls GDT for invoice status.
func (c *Client) GetInvoiceStatus(ctx context.Context, transactionID string) (*domain.GDTStatusResponse, error) {
	var out domain.GDTStatusResponse
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

// SubmitDeclaration submits a signed declaration XML (HTKK file) to GDT.
func (c *Client) SubmitDeclaration(ctx context.Context, declarationXML, certID string) (*domain.GDTDeclarationSubmitResponse, error) {
	body, _ := json.Marshal(SubmitRequest{XML: declarationXML, CertID: certID})
	var out domain.GDTDeclarationSubmitResponse
	if err := c.do(ctx, http.MethodPost, "/api/submission/declare", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// QueryDeclarationStatus polls GDT for declaration status.
func (c *Client) QueryDeclarationStatus(ctx context.Context, submissionID string) (*domain.GDTDeclarationStatusResponse, error) {
	var out domain.GDTDeclarationStatusResponse
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
		req.Header.Set("X-Request-ID", fmt.Sprintf("gotax-%d", time.Now().UnixNano()))
		if c.authToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.authToken)
		}
		if c.logger != nil {
			c.logger("gdt: %s %s attempt=%d/%d", method, path, i+1, attempts)
		}
		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = domain.ErrGDTUnavailable // network error → retry
			if c.logger != nil {
				c.logger("gdt: %s %s error=%v (retry)", method, path, err)
			}
		} else {
			respBody, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			bodyStr := string(respBody)
			if len(bodyStr) > 4096 {
				bodyStr = bodyStr[:4096] + "...(truncated)"
			}
			switch {
			case resp.StatusCode >= 200 && resp.StatusCode < 300:
				if readErr != nil {
					return fmt.Errorf("gdt: read response: %w", readErr)
				}
				if out != nil && len(respBody) > 0 {
					if err := json.Unmarshal(respBody, out); err != nil {
						return &GDTError{Err: fmt.Errorf("gdt: decode response: %w", err), StatusCode: resp.StatusCode, Body: bodyStr, Path: path}
					}
				}
				if c.logger != nil {
					c.logger("gdt: %s %s → %d OK", method, path, resp.StatusCode)
				}
				return nil
			case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
				return &GDTError{Err: domain.ErrGDTUnauthorized, StatusCode: resp.StatusCode, Body: bodyStr, Path: path}
			case resp.StatusCode >= 400 && resp.StatusCode < 500:
				return &GDTError{Err: domain.ErrGDTRejected, StatusCode: resp.StatusCode, Body: bodyStr, Path: path}
			default: // 5xx → retry
				lastErr = &GDTError{Err: domain.ErrGDTUnavailable, StatusCode: resp.StatusCode, Body: bodyStr, Path: path}
				if c.logger != nil {
					c.logger("gdt: %s %s → %d (retry)", method, path, resp.StatusCode)
				}
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
