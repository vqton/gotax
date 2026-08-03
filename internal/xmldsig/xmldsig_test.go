package xmldsig

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

const txml = `<?xml version="1.0" encoding="UTF-8"?>
<BK:HoaDon xmlns:BK="http://gdt.gov.vn/HoaDon">
  <BK:ThongTinChung>
    <BK:MaSoThue>0123456789</BK:MaSoThue>
    <BK:MauSo>01GTKT0/001</BK:MauSo>
  </BK:ThongTinChung>
</BK:HoaDon>`

func TestSignVerifyRoundTrip(t *testing.T) {
	key := testKey(t)
	canon, err := Canonicalize([]byte(txml))
	require.NoError(t, err)

	sig, err := Sign(key, canon)
	require.NoError(t, err)
	assert.NotEmpty(t, sig)

	require.NoError(t, Verify(&key.PublicKey, canon, sig))
}

func TestVerifyRejectsTampered(t *testing.T) {
	key := testKey(t)
	canon, _ := Canonicalize([]byte(txml))
	sig, _ := Sign(key, canon)

	tampered := append([]byte(nil), canon...)
	tampered[bytes.Index(tampered, []byte("0123456789"))] = '9'
	assert.Error(t, Verify(&key.PublicKey, tampered, sig))

	// wrong key
	other := testKey(t)
	assert.Error(t, Verify(&other.PublicKey, canon, sig))
}

func TestSignDeterministic(t *testing.T) {
	key := testKey(t)
	canon, _ := Canonicalize([]byte(txml))
	s1, _ := Sign(key, canon)
	s2, _ := Sign(key, canon)
	assert.Equal(t, s1, s2) // PKCS1v15 is deterministic
}

func TestCanonicalizeWhitespaceInsensitive(t *testing.T) {
	// inter-element whitespace is significant in C14N — preserved. Same
	// input must canonicalize deterministically (stable output).
	compact := `<?xml version="1.0"?><BK:HoaDon xmlns:BK="http://gdt.gov.vn/HoaDon"><BK:ThongTinChung><BK:MaSoThue>0123456789</BK:MaSoThue></BK:ThongTinChung></BK:HoaDon>`
	c1, err := Canonicalize([]byte(txml))
	require.NoError(t, err)
	c2, err := Canonicalize([]byte(txml))
	require.NoError(t, err)
	assert.Equal(t, string(c1), string(c2))
	// pretty input keeps its inter-element whitespace; prefixes preserved
	assert.Contains(t, string(c1), "<BK:HoaDon xmlns:BK=\"http://gdt.gov.vn/HoaDon\">\n  <BK:ThongTinChung>\n    <BK:MaSoThue>")
	// compact input stays compact
	c3, err := Canonicalize([]byte(compact))
	require.NoError(t, err)
	assert.Equal(t, "<BK:HoaDon xmlns:BK=\"http://gdt.gov.vn/HoaDon\"><BK:ThongTinChung><BK:MaSoThue>0123456789</BK:MaSoThue></BK:ThongTinChung></BK:HoaDon>", string(c3))
}

func TestCanonicalizeNormalizes(t *testing.T) {
	input := "\xEF\xBB\xBF<Root xmlns=\"urn:test\" b=\"2\" a=\"1\">\r\ntext &amp; more\r</Root>"
	canon, err := Canonicalize([]byte(input))
	require.NoError(t, err)
	out := string(canon)
	// BOM stripped, CRLF → LF
	assert.NotContains(t, out, "\xEF\xBB\xBF")
	assert.NotContains(t, out, "\r")
	// default namespace on root: bare element name, xmlns kept
	assert.Contains(t, out, `<Root xmlns="urn:test" a="1" b="2">`)
	// text re-escaped, LF preserved
	assert.Contains(t, out, "text &amp; more\n")
}

func TestCanonicalizeEscapesText(t *testing.T) {
	// decoder yields raw chars → re-escaped in canonical output
	canon, err := Canonicalize([]byte(`<Doc>a&lt;b &amp; c &gt; d</Doc>`))
	require.NoError(t, err)
	assert.Equal(t, `<Doc>a&lt;b &amp; c &gt; d</Doc>`, string(canon))
}

func TestCanonicalizeEmptyElement(t *testing.T) {
	canon, err := Canonicalize([]byte(`<Doc><Empty/></Doc>`))
	require.NoError(t, err)
	assert.Equal(t, `<Doc><Empty></Empty></Doc>`, string(canon))
}

func TestCanonicalizeRejectsMalformed(t *testing.T) {
	_, err := Canonicalize([]byte(`<Root><Unclosed>`))
	assert.Error(t, err)
}

func TestDigest(t *testing.T) {
	key := testKey(t)
	canon, _ := Canonicalize([]byte(txml))
	sig, err := Sign(key, canon)
	require.NoError(t, err)
	// SHA-256 digest is 32 bytes
	d := Digest(canon)
	assert.Len(t, d, 32)
	_ = sig
}

func TestParsePrivateKeyPEM(t *testing.T) {
	key := testKey(t)
	pkcs1 := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	got, err := ParsePrivateKeyPEM(pkcs1)
	require.NoError(t, err)
	canon, _ := Canonicalize([]byte(txml))
	sig, err := Sign(got, canon)
	require.NoError(t, err)
	require.NoError(t, Verify(&key.PublicKey, canon, sig))

	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	got8, err := ParsePrivateKeyPEM(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8}))
	require.NoError(t, err)
	sig8, err := Sign(got8, canon)
	require.NoError(t, err)
	require.NoError(t, Verify(&key.PublicKey, canon, sig8))

	_, err = ParsePrivateKeyPEM([]byte("not pem"))
	assert.Error(t, err)
}
