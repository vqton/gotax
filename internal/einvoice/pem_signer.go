package einvoice

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"gotax/internal/xmldsig"
)

// SignResult is a detached signature over a canonicalized document.
type SignResult struct {
	SignatureBase64 string
	SignedAt        string
}

// pemSigner signs canonicalized XML with an RSA private key. SignTXML embeds
// the signature into BK:ChuKySo before BK:KyThuat closes (TAX_TEMPLATES §5);
// SignDocument returns the detached base64 signature (HTKK declarations).
type pemSigner struct {
	key    *rsa.PrivateKey
	serial string
	now    func() time.Time
}

// NewPEMSigner builds a signer from an RSA private key. serial is the
// digital certificate serial placed in BK:SerialNumber.
func NewPEMSigner(key *rsa.PrivateKey, serial string, now func() time.Time) *pemSigner {
	return &pemSigner{key: key, serial: serial, now: now}
}

func (p *pemSigner) SignDocument(xmlBody string) (SignResult, error) {
	canon, err := xmldsig.Canonicalize([]byte(xmlBody))
	if err != nil {
		return SignResult{}, fmt.Errorf("sign document: %w", err)
	}
	sig, err := xmldsig.SignBase64(p.key, canon)
	if err != nil {
		return SignResult{}, fmt.Errorf("sign document: %w", err)
	}
	return SignResult{
		SignatureBase64: sig,
		SignedAt:        p.now().Format(time.RFC3339),
	}, nil
}

func (p *pemSigner) SignTXML(xmlBody, signatureID string) (string, error) {
	canon, err := xmldsig.Canonicalize([]byte(xmlBody))
	if err != nil {
		return "", fmt.Errorf("sign TXML: %w", err)
	}
	sig, err := xmldsig.SignBase64(p.key, canon)
	if err != nil {
		return "", fmt.Errorf("sign TXML: %w", err)
	}
	stamp := p.now().Format("2006-01-02T15:04:05-07:00")
	chuku := "<BK:ChuKySo><BK:SerialNumber>" + p.serial +
		"</BK:SerialNumber><BK:ThoiDiemKy>" + stamp +
		"</BK:ThoiDiemKy><BK:DuLieuKy>" + sig +
		"</BK:DuLieuKy></BK:ChuKySo>"
	signed := strings.Replace(string(canon), "</BK:KyThuat>", chuku+"</BK:KyThuat>", 1)
	if signed == string(canon) {
		return "", fmt.Errorf("sign TXML: BK:KyThuat not found")
	}
	return signed, nil
}
