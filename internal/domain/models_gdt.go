package domain

// GDT wire DTOs (TAX_SPECS §4). Live in domain so the service layer can
// depend on them without importing the net/http transport adapter.
// Response codes per spec §4.2: 00 acknowledged, 01 schema failure,
// 02 duplicate, 03 tax code not found, 10 period already declared,
// 99 system error.

// GDTSubmitResponse is the GDT acknowledgement of an e-invoice submission.
type GDTSubmitResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	GDTRef        string `json:"gdt_ref"`
	Message       string `json:"message,omitempty"`
}

// GDTStatusResponse is the e-invoice status payload.
type GDTStatusResponse struct {
	Status  string `json:"status"`
	GDTRef  string `json:"gdt_ref,omitempty"`
	Message string `json:"message,omitempty"`
}

// GDTDeclarationSubmitResponse is the GDT acknowledgement of a declaration
// submission.
type GDTDeclarationSubmitResponse struct {
	SubmissionID string `json:"submission_id"`
	Status       string `json:"status"`
	Code         string `json:"code,omitempty"`
	AckRef       string `json:"ack_ref,omitempty"`
	Message      string `json:"message,omitempty"`
}

// GDTDeclarationStatusResponse is the declaration status payload.
type GDTDeclarationStatusResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code,omitempty"`
	AckRef  string `json:"ack_ref,omitempty"`
	Message string `json:"message,omitempty"`
}

// TaxCodeLookupResponse is the GDT tax code validation result.
// Per GDT spec: status 00=active, 01=not found, 02=inactive.
type TaxCodeLookupResponse struct {
	TaxCode     string `json:"tax_code"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status"`      // "ACTIVE", "INACTIVE", "NOT_FOUND"
	Message     string `json:"message,omitempty"`
	LastUpdated string `json:"last_updated,omitempty"`
}
