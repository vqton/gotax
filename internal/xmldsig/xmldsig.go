// Package xmldsig provides XML canonicalization (exclusive C14N 1.0) and
// RSA-SHA256 signing primitives used for GDT e-invoice TXML signatures
// (Decree 254/2026/ND-CP, BK:ChuKySo/DuLieuKy = base64 signature).
package xmldsig

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
)

const xmlNamespace = "http://www.w3.org/XML/1998/namespace"

// Canonicalize applies exclusive XML canonicalization 1.0: BOM stripped,
// CR/CRLF → LF, CDATA → text, attributes sorted by (namespace URI, local
// name), namespace declarations rendered only where visible (inherited
// prefixes rendered at first use), text re-escaped. Deterministic output —
// signing the canonical form is reproducible.
func Canonicalize(data []byte) ([]byte, error) {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.Strict = true
	dec.CharsetReader = func(cs string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(cs, "utf-8") || strings.EqualFold(cs, "utf8") {
			return input, nil
		}
		return nil, fmt.Errorf("xmldsig: unsupported charset %q", cs)
	}
	sc := &canonScope{rendered: map[string]bool{}}
	var buf bytes.Buffer
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("xmldsig: parse: %w", err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if err := renderStart(&buf, sc, t); err != nil {
				return nil, err
			}
		case xml.EndElement:
			prefix := sc.prefixFor(t.Name.Space)
			buf.WriteString("</")
			if prefix != "" {
				buf.WriteString(prefix)
				buf.WriteByte(':')
			}
			buf.WriteString(t.Name.Local)
			buf.WriteByte('>')
			sc.pop()
		case xml.CharData:
			norm := bytes.ReplaceAll(t, []byte{'\r', '\n'}, []byte{'\n'})
			norm = bytes.ReplaceAll(norm, []byte{'\r'}, []byte{'\n'})
			xmlEscape(&buf, string(norm))
		case xml.Comment:
			buf.WriteString("<!--")
			buf.Write(t)
			buf.WriteString("-->")
		case xml.ProcInst:
			if strings.EqualFold(t.Target, "xml") {
				continue // drop XML declaration
			}
			buf.WriteString("<?")
			buf.WriteString(t.Target)
			buf.WriteByte(' ')
			buf.Write(t.Inst)
			buf.WriteString("?>")
		case xml.Directive:
			continue // drop DOCTYPE/DTD
		}
	}
	return buf.Bytes(), nil
}

type nsDecl struct{ prefix, uri string }

type canonScope struct {
	stack    []map[string]string // prefix → uri, one level per element depth
	rendered map[string]bool     // "prefix\x00uri" already emitted in document
}

func (s *canonScope) push(bindings map[string]string) { s.stack = append(s.stack, bindings) }
func (s *canonScope) pop()                            { s.stack = s.stack[:len(s.stack)-1] }

// prefixFor returns the in-scope prefix bound to uri ("" if default/unbound).
// Deterministic: lexicographically smallest prefix among bindings.
func (s *canonScope) prefixFor(uri string) string {
	if uri == "" {
		return ""
	}
	best := ""
	for i := len(s.stack) - 1; i >= 0; i-- {
		for p, u := range s.stack[i] {
			if u == uri && (best == "" || p < best) {
				best = p
			}
		}
	}
	return best
}

func (s *canonScope) ancestorURI(prefix string) (string, bool) {
	for i := len(s.stack) - 1; i >= 0; i-- {
		if uri, ok := s.stack[i][prefix]; ok {
			return uri, true
		}
	}
	return "", false
}

func renderStart(buf *bytes.Buffer, sc *canonScope, t xml.StartElement) error {
	// separate xmlns bindings declared on this element (source order)
	bindings := []nsDecl{}
	attrs := []xml.Attr{}
	for _, a := range t.Attr {
		switch {
		case a.Name.Space == "xmlns":
			bindings = append(bindings, nsDecl{prefix: a.Name.Local, uri: a.Value})
		case a.Name.Space == "" && a.Name.Local == "xmlns":
			bindings = append(bindings, nsDecl{prefix: "", uri: a.Value})
		default:
			attrs = append(attrs, a)
		}
	}
	bindMap := map[string]string{}
	for _, b := range bindings {
		bindMap[b.prefix] = b.uri
	}
	sc.push(bindMap)

	// namespaces to render: (a) bindings that change scope, (b) inherited
	// prefixes first used by this element.
	toRender := []nsDecl{}
	seen := map[string]bool{}
	add := func(prefix, uri string) {
		k := prefix + "\x00" + uri
		if seen[k] {
			return
		}
		seen[k] = true
		toRender = append(toRender, nsDecl{prefix, uri})
	}
	for _, b := range bindings {
		if anc, ok := sc.ancestorURI(b.prefix); !ok || anc != b.uri {
			add(b.prefix, b.uri)
		}
	}
	used := []nsDecl{}
	if t.Name.Space != "" {
		used = append(used, nsDecl{sc.prefixFor(t.Name.Space), t.Name.Space})
	}
	for _, a := range attrs {
		if a.Name.Space != "" && a.Name.Space != xmlNamespace {
			used = append(used, nsDecl{sc.prefixFor(a.Name.Space), a.Name.Space})
		}
	}
	for _, u := range used {
		k := u.prefix + "\x00" + u.uri
		if u.uri != "" && u.uri != xmlNamespace && !sc.rendered[k] {
			add(u.prefix, u.uri)
			sc.rendered[k] = true
		}
	}
	sort.SliceStable(toRender, func(i, j int) bool { return toRender[i].prefix < toRender[j].prefix })

	elPrefix := sc.prefixFor(t.Name.Space)
	buf.WriteByte('<')
	if elPrefix != "" {
		buf.WriteString(elPrefix)
		buf.WriteByte(':')
	}
	buf.WriteString(t.Name.Local)
	for _, b := range toRender {
		buf.WriteString(" xmlns")
		if b.prefix != "" {
			buf.WriteByte(':')
			buf.WriteString(b.prefix)
		}
		buf.WriteString(`="`)
		xmlEscape(buf, b.uri)
		buf.WriteByte('"')
	}
	sort.SliceStable(attrs, func(i, j int) bool {
		ni, nj := attrs[i].Name, attrs[j].Name
		if ni.Space != nj.Space {
			return ni.Space < nj.Space
		}
		return ni.Local < nj.Local
	})
	for _, a := range attrs {
		buf.WriteByte(' ')
		if a.Name.Space != "" {
			p := sc.prefixFor(a.Name.Space)
			if p != "" {
				buf.WriteString(p)
				buf.WriteByte(':')
			}
		}
		buf.WriteString(a.Name.Local)
		buf.WriteString(`="`)
		xmlEscape(buf, a.Value)
		buf.WriteByte('"')
	}
	buf.WriteByte('>')
	return nil
}

func xmlEscape(buf *bytes.Buffer, s string) {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "]]>", "]]&gt;")
	buf.WriteString(s)
}

// Digest returns the SHA-256 digest of data.
func Digest(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// Sign returns the RSA-SHA256 (PKCS#1 v1.5) signature over data.
func Sign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	digest := Digest(data)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
}

// Verify checks an RSA-SHA256 signature over data.
func Verify(pub *rsa.PublicKey, data, sig []byte) error {
	digest := Digest(data)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest, sig)
}

// SignBase64 returns the signature as base64 (BK:DuLieuKy payload).
func SignBase64(key *rsa.PrivateKey, data []byte) (string, error) {
	sig, err := Sign(key, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// VerifyBase64 verifies a base64-encoded signature.
func VerifyBase64(pub *rsa.PublicKey, data []byte, sigB64 string) error {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return fmt.Errorf("xmldsig: base64: %w", err)
	}
	return Verify(pub, data, sig)
}

// ParsePrivateKeyPEM parses an RSA private key from PKCS#1 or PKCS#8 PEM.
func ParsePrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("xmldsig: no PEM block found")
	}
	if k, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rk, ok := k.(*rsa.PrivateKey); ok {
			return rk, nil
		}
		return nil, fmt.Errorf("xmldsig: not an RSA key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
