package accesstoken

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type Header struct {
	Ver string `json:"ver"`
	Typ string `json:"typ"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

type Payload struct {
	Jti string      `json:"jti"`
	Sub string      `json:"sub"`
	Aud interface{} `json:"aud"`
	Scp string      `json:"scp"`
	Iss string      `json:"iss"`
	Iat float64     `json:"iat"`
	Exp float64     `json:"exp"`
	Ext string      `json:"ext,omitempty"`
}

type ParsedExt struct {
	Revocable  interface{} `json:"revocable,omitempty"`
	Sid        string      `json:"sid,omitempty"`
	Ssi        string      `json:"ssi,omitempty"`
	Rpp        interface{} `json:"rpp,omitempty"`
	Rlm        string      `json:"rlm,omitempty"`
	Llt        interface{} `json:"llt,omitempty"`
	Pud        interface{} `json:"pud,omitempty"`
	ProjectKey string      `json:"projectKey,omitempty"`
}

type ParsedToken struct {
	Header    Header
	Payload   Payload
	ParsedExt *ParsedExt
	Claims    map[string]interface{}
}

func (p *Payload) AudAsStringSlice() []string {
	switch v := p.Aud.(type) {
	case string:
		return []string{v}
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	default:
		return nil
	}
}

func (p *Payload) AudAsString() string {
	return strings.Join(p.AudAsStringSlice(), " ")
}

func (p *Payload) IatTime() time.Time {
	return time.Unix(int64(p.Iat), 0)
}

func (p *Payload) ExpTime() time.Time {
	return time.Unix(int64(p.Exp), 0)
}

func (p *Payload) Scopes() []string {
	if p.Scp == "" {
		return nil
	}
	return strings.Fields(p.Scp)
}

// Not jwt.DecodeSegment: its behaviour depends on the package-level DecodePaddingAllowed and DecodeStrict settings.
func DecodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(strings.TrimRight(seg, "="))
}

func ParseToken(rawToken string) (*ParsedToken, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return nil, errorutils.CheckErrorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	headerBytes, err := DecodeSegment(parts[0])
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to decode JWT header: %w", err)
	}
	var header Header
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse JWT header: %w", err)
	}

	payloadBytes, err := DecodeSegment(parts[1])
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to decode JWT payload: %w", err)
	}
	var payload Payload
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errorutils.CheckErrorf("failed to parse JWT payload: %w", err)
	}
	claims, err := decodeClaims(payloadBytes)
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to parse JWT payload: %w", err)
	}

	parsed := &ParsedToken{Header: header, Payload: payload, Claims: claims}
	if payload.Ext != "" {
		var ext ParsedExt
		if err = json.Unmarshal([]byte(payload.Ext), &ext); err != nil {
			return nil, errorutils.CheckErrorf("failed to parse ext field as JSON: %w", err)
		}
		parsed.ParsedExt = &ext
		extClaims, err := decodeClaims([]byte(payload.Ext))
		if err != nil {
			return nil, errorutils.CheckErrorf("failed to parse ext field as JSON: %w", err)
		}
		claims["ext"] = extClaims
	}
	return parsed, nil
}

func decodeClaims(data []byte) (map[string]interface{}, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	claims := map[string]interface{}{}
	if err := decoder.Decode(&claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func ValidateHeader(header Header) error {
	if header.Alg != "RS256" {
		return errorutils.CheckErrorf("invalid header: expected alg \"RS256\", got %q", header.Alg)
	}
	if header.Typ != "JWT" {
		return errorutils.CheckErrorf("invalid header: expected typ \"JWT\", got %q", header.Typ)
	}
	if header.Ver != "2" {
		return errorutils.CheckErrorf("invalid header: expected ver \"2\", got %q", header.Ver)
	}
	return nil
}

func ValidateSignature(rawToken string, key *rsa.PublicKey) error {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}), jwt.WithoutClaimsValidation())
	if _, err := parser.Parse(rawToken, func(*jwt.Token) (interface{}, error) { return key, nil }); err != nil {
		return errorutils.CheckErrorf("invalid signature: %w", err)
	}
	return nil
}

func ParseRootCertificate(rootCertBase64 string) (*rsa.PublicKey, error) {
	rootCertRaw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(rootCertBase64))
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to base64-decode root certificate: %w", err)
	}
	certificate, err := x509.ParseCertificate(rootCertRaw)
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to parse X.509 certificate: %w", err)
	}
	rsaKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errorutils.CheckErrorf("root certificate public key is not RSA")
	}
	return rsaKey, nil
}

func ValidatePayloadFields(payload Payload) error {
	switch {
	case payload.Jti == "":
		return errorutils.CheckErrorf("missing required field: jti")
	case payload.Sub == "":
		return errorutils.CheckErrorf("missing required field: sub")
	case payload.AudAsString() == "":
		return errorutils.CheckErrorf("missing required field: aud")
	case payload.Scp == "":
		return errorutils.CheckErrorf("missing required field: scp (empty scope)")
	case payload.Iss == "":
		return errorutils.CheckErrorf("missing required field: iss")
	case payload.Iat == 0:
		return errorutils.CheckErrorf("missing required field: iat")
	case payload.Exp != 0 && payload.Exp <= payload.Iat:
		return errorutils.CheckErrorf("invalid token: exp (%d) must be greater than iat (%d)", int64(payload.Exp), int64(payload.Iat))
	}
	return nil
}

func ValidateExpiration(payload Payload, now time.Time) error {
	if payload.Exp == 0 {
		return nil
	}
	if int64(payload.Exp) <= now.Unix() {
		return errorutils.CheckErrorf("token expired at %s", payload.ExpTime().UTC().Format(time.RFC3339))
	}
	return nil
}
