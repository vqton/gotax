package gdt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gotax/internal/domain"
)

// mockServer implements the GDT wire contract with configurable failures.
type mockServer struct {
	t            *testing.T
	failSubmits  atomic.Int32 // 502 responses before success
	submitCalls  atomic.Int32
	submitBody   atomic.Value
	statusCalls  atomic.Int32
	cancelCalls  atomic.Int32
	declSubmitCalls atomic.Int32
	declSubmitBody  atomic.Value
	declStatusCalls atomic.Int32
	authRequired bool
	server       *httptest.Server
}

func newMockServer(t *testing.T) *mockServer {
	m := &mockServer{t: t}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/invoice/submit", func(w http.ResponseWriter, r *http.Request) {
		m.submitCalls.Add(1)
		if m.authRequired && r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var req SubmitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.submitBody.Store(req.XML)
		if m.failSubmits.Load() > 0 {
			m.failSubmits.Add(-1)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		resp := domain.GDTSubmitResponse{TransactionID: "TXN-1", Status: "SUBMITTED", GDTRef: "GDT-REF-1"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/api/invoice/status", func(w http.ResponseWriter, r *http.Request) {
		m.statusCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.GDTStatusResponse{Status: "ACKNOWLEDGED", GDTRef: "GDT-REF-1"})
	})
	mux.HandleFunc("/api/invoice/cancel", func(w http.ResponseWriter, r *http.Request) {
		m.cancelCalls.Add(1)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/api/submission/declare", func(w http.ResponseWriter, r *http.Request) {
		m.declSubmitCalls.Add(1)
		var req SubmitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.declSubmitBody.Store(req.XML)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.GDTDeclarationSubmitResponse{SubmissionID: "SUB-1", Status: "SUBMITTED", Code: "00", AckRef: "ACK-REF-1"})
	})
	mux.HandleFunc("/api/submission/status", func(w http.ResponseWriter, r *http.Request) {
		m.declStatusCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.GDTDeclarationStatusResponse{Status: "ACKNOWLEDGED", Code: "00", AckRef: "ACK-REF-1"})
	})
	m.server = httptest.NewServer(mux)
	return m
}

func (m *mockServer) URL() string { return m.server.URL }
func (m *mockServer) Close()      { m.server.Close() }

func testClient(t *testing.T, url string, opts ...Option) *Client {
	c, err := New(url, append(opts, WithRetry(time.Millisecond, 2*time.Millisecond))...)
	require.NoError(t, err)
	return c
}

func TestSubmitInvoice(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	c := testClient(t, m.URL())

	resp, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	require.NoError(t, err)
	assert.Equal(t, "TXN-1", resp.TransactionID)
	assert.Equal(t, "SUBMITTED", resp.Status)
	assert.Equal(t, "GDT-REF-1", resp.GDTRef)
	assert.Equal(t, "<invoice/>", m.submitBody.Load().(string))
	assert.Equal(t, int32(1), m.submitCalls.Load())
}

func TestSubmitRetriesOn5xx(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	m.failSubmits.Store(2)
	c := testClient(t, m.URL())

	resp, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	require.NoError(t, err)
	assert.Equal(t, "TXN-1", resp.TransactionID)
	assert.Equal(t, int32(3), m.submitCalls.Load()) // 2 failures + 1 success
}

func TestSubmitRetriesExhausted(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	m.failSubmits.Store(99)
	c := testClient(t, m.URL())

	_, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	assert.ErrorIs(t, err, domain.ErrGDTUnavailable)
	assert.Equal(t, int32(3), m.submitCalls.Load()) // max 3 attempts
}

func TestSubmitUnauthorized(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	m.authRequired = true
	c := testClient(t, m.URL())

	_, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	assert.ErrorIs(t, err, domain.ErrGDTUnauthorized)
	assert.Equal(t, int32(1), m.submitCalls.Load()) // no retry on 401
}

func TestSubmitBadRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()
	c := testClient(t, srv.URL)

	_, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	assert.ErrorIs(t, err, domain.ErrGDTRejected)
}

func TestSubmitMalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()
	c := testClient(t, srv.URL)

	_, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrGDTUnavailable)
}

func TestSubmitNetworkError(t *testing.T) {
	// closed server → connection refused; retried then domain.ErrGDTUnavailable
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()
	c := testClient(t, url)

	_, err := c.SubmitInvoice(context.Background(), "<invoice/>", "sig-1")
	assert.ErrorIs(t, err, domain.ErrGDTUnavailable)
}

func TestSubmitContextCancel(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	m.failSubmits.Store(99)
	c := testClient(t, m.URL())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.SubmitInvoice(ctx, "<invoice/>", "sig-1")
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGetInvoiceStatus(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	c := testClient(t, m.URL())

	st, err := c.GetInvoiceStatus(context.Background(), "TXN-1")
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", st.Status)
	assert.Equal(t, int32(1), m.statusCalls.Load())
}

func TestCancelInvoice(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	c := testClient(t, m.URL())

	require.NoError(t, c.CancelInvoice(context.Background(), "TXN-1", "buyer request"))
	assert.Equal(t, int32(1), m.cancelCalls.Load())
}

func TestSubmitDeclaration(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	c := testClient(t, m.URL())

	resp, err := c.SubmitDeclaration(context.Background(), "<BoKe/>", "sig-1")
	require.NoError(t, err)
	assert.Equal(t, "SUB-1", resp.SubmissionID)
	assert.Equal(t, "SUBMITTED", resp.Status)
	assert.Equal(t, "00", resp.Code)
	assert.Equal(t, "<BoKe/>", m.declSubmitBody.Load().(string))
	assert.Equal(t, int32(1), m.declSubmitCalls.Load())
}

func TestQueryDeclarationStatus(t *testing.T) {
	m := newMockServer(t)
	defer m.Close()
	c := testClient(t, m.URL())

	st, err := c.QueryDeclarationStatus(context.Background(), "SUB-1")
	require.NoError(t, err)
	assert.Equal(t, "ACKNOWLEDGED", st.Status)
	assert.Equal(t, "00", st.Code)
	assert.Equal(t, "ACK-REF-1", st.AckRef)
	assert.Equal(t, int32(1), m.declStatusCalls.Load())
}

func TestNewInvalidURL(t *testing.T) {
	_, err := New("://bad", WithRetry(time.Millisecond))
	assert.Error(t, err)
}
