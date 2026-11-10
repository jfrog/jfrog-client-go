package accesstoken

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestCert(t *testing.T) (string, *rsa.PrivateKey) {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	certDER, err := x509.CreateCertificate(rand.Reader, certTemplate(), certTemplate(), &privKey.PublicKey, privKey)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(certDER), privKey
}

func certTemplate() *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}
}

func publicKey(t *testing.T, certB64 string) *rsa.PublicKey {
	t.Helper()
	key, err := ParseRootCertificate(certB64)
	require.NoError(t, err)
	return key
}

func createSignedToken(t *testing.T, privKey *rsa.PrivateKey, header map[string]interface{}, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	for k, v := range header {
		token.Header[k] = v
	}
	tokenStr, err := token.SignedString(privKey)
	require.NoError(t, err)
	return tokenStr
}

func standardHeader() map[string]interface{} {
	return map[string]interface{}{"ver": "2", "kid": "test-kid"}
}

func standardClaims() jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"jti": "test-jti-12345",
		"sub": "jfac@01ftqhpvg3kn760kg8qbt411qy/users/admin",
		"aud": "*@*",
		"scp": "applied-permissions/admin",
		"iss": "jfac@01ftqhpvg3kn760kg8qbt411qy",
		"iat": float64(now.Unix()),
		"exp": float64(now.Add(1 * time.Hour).Unix()),
		"ext": `{"revocable":"true"}`,
	}
}

func encodeSegment(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func TestParseToken_ValidToken(t *testing.T) {
	_, privKey := generateTestCert(t)
	parsed, err := ParseToken(createSignedToken(t, privKey, standardHeader(), standardClaims()))
	require.NoError(t, err)
	assert.Equal(t, "RS256", parsed.Header.Alg)
	assert.Equal(t, "JWT", parsed.Header.Typ)
	assert.Equal(t, "2", parsed.Header.Ver)
	assert.Equal(t, "test-jti-12345", parsed.Payload.Jti)
	assert.Equal(t, "applied-permissions/admin", parsed.Payload.Scp)
	require.NotNil(t, parsed.ParsedExt)
}

func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := ParseToken("not-a-jwt")
	assert.Error(t, err)
}

func TestParseToken_InvalidExtJSON(t *testing.T) {
	_, privKey := generateTestCert(t)
	claims := standardClaims()
	claims["ext"] = "not-valid-json"
	_, err := ParseToken(createSignedToken(t, privKey, standardHeader(), claims))
	assert.Error(t, err)
}

func TestParseToken_NoExt(t *testing.T) {
	_, privKey := generateTestCert(t)
	claims := standardClaims()
	delete(claims, "ext")
	parsed, err := ParseToken(createSignedToken(t, privKey, standardHeader(), claims))
	require.NoError(t, err)
	assert.Nil(t, parsed.ParsedExt)
}

func TestParseToken_Base64URLSegments(t *testing.T) {
	segment := encodeSegment(`{"sub":"jfac@01h/users/דנה","scp":"applied-permissions/user"}`)
	require.True(t, strings.ContainsAny(segment, "-_"))
	parsed, err := ParseToken("e30." + segment + ".sig")
	require.NoError(t, err)
	assert.Equal(t, "jfac@01h/users/דנה", parsed.Payload.Sub)
}

func TestParseToken_KeepsUnknownClaims(t *testing.T) {
	_, privKey := generateTestCert(t)
	claims := standardClaims()
	claims["foo"] = "bar"
	parsed, err := ParseToken(createSignedToken(t, privKey, standardHeader(), claims))
	require.NoError(t, err)
	assert.Equal(t, "bar", parsed.Claims["foo"])
	assert.IsType(t, map[string]interface{}{}, parsed.Claims["ext"])
}

func TestDecodeSegment(t *testing.T) {
	for _, tc := range []struct {
		name    string
		segment string
		want    []byte
	}{
		{"unpadded", "eyJhIjoxfQ", []byte(`{"a":1}`)},
		{"padded", "eyJhIjoxfQ==", []byte(`{"a":1}`)},
		{"url alphabet", "-_8", []byte{0xfb, 0xff}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodeSegment(tc.segment)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
	_, err := DecodeSegment("+/8")
	assert.Error(t, err)
}

func TestValidateHeader_Valid(t *testing.T) {
	assert.NoError(t, ValidateHeader(Header{Ver: "2", Typ: "JWT", Alg: "RS256", Kid: "test"}))
}

func TestValidateHeader_WrongAlg(t *testing.T) {
	assert.Error(t, ValidateHeader(Header{Ver: "2", Typ: "JWT", Alg: "HS256", Kid: "test"}))
}

func TestValidateHeader_WrongTyp(t *testing.T) {
	assert.Error(t, ValidateHeader(Header{Ver: "2", Typ: "JWS", Alg: "RS256", Kid: "test"}))
}

func TestValidateHeader_WrongVer(t *testing.T) {
	assert.ErrorContains(t, ValidateHeader(Header{Ver: "1", Typ: "JWT", Alg: "RS256", Kid: "test"}), `expected ver "2"`)
}

func TestValidateHeader_MissingFields(t *testing.T) {
	assert.Error(t, ValidateHeader(Header{}))
}

func TestValidateSignature_Valid(t *testing.T) {
	certB64, privKey := generateTestCert(t)
	token := createSignedToken(t, privKey, standardHeader(), standardClaims())
	assert.NoError(t, ValidateSignature(token, publicKey(t, certB64)))
}

func TestValidateSignature_WrongCert(t *testing.T) {
	_, privKey := generateTestCert(t)
	otherCertB64, _ := generateTestCert(t)
	token := createSignedToken(t, privKey, standardHeader(), standardClaims())
	assert.ErrorContains(t, ValidateSignature(token, publicKey(t, otherCertB64)), "invalid signature")
}

func TestValidateSignature_RejectsAlgNone(t *testing.T) {
	certB64, _ := generateTestCert(t)
	token := encodeSegment(`{"alg":"none","typ":"JWT","ver":"2"}`) + "." + encodeSegment(`{"sub":"x"}`) + "."
	assert.ErrorContains(t, ValidateSignature(token, publicKey(t, certB64)), "invalid signature")
}

func TestValidateSignature_RejectsHS256(t *testing.T) {
	certB64, _ := generateTestCert(t)
	certDER, err := base64.StdEncoding.DecodeString(certB64)
	require.NoError(t, err)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims()).SignedString(certDER)
	require.NoError(t, err)
	assert.ErrorContains(t, ValidateSignature(token, publicKey(t, certB64)), "invalid signature")
}

func TestValidatePayloadFields_Valid(t *testing.T) {
	p := Payload{
		Jti: "jti-1",
		Sub: "jfac@01ftqhpvg3kn760kg8qbt411qy/users/admin",
		Aud: "*@*",
		Scp: "applied-permissions/admin",
		Iss: "jfac@01ftqhpvg3kn760kg8qbt411qy",
		Iat: float64(time.Now().Unix()),
		Exp: float64(time.Now().Add(1 * time.Hour).Unix()),
	}
	assert.NoError(t, ValidatePayloadFields(p))
}

func TestValidatePayloadFields_MissingJti(t *testing.T) {
	p := Payload{Sub: "sub", Aud: "aud", Scp: "scp", Iss: "iss", Iat: 1, Exp: 2}
	assert.EqualError(t, ValidatePayloadFields(p), "missing required field: jti")
}

func TestValidatePayloadFields_MissingSub(t *testing.T) {
	p := Payload{Jti: "jti", Aud: "aud", Scp: "scp", Iss: "iss", Iat: 1, Exp: 2}
	assert.EqualError(t, ValidatePayloadFields(p), "missing required field: sub")
}

func TestValidatePayloadFields_EmptyScp(t *testing.T) {
	p := Payload{Jti: "jti", Sub: "sub", Aud: "aud", Iss: "iss", Iat: 1, Exp: 2}
	assert.EqualError(t, ValidatePayloadFields(p), "missing required field: scp (empty scope)")
}

func TestValidatePayloadFields_ExpNotGreaterThanIat(t *testing.T) {
	p := Payload{Jti: "jti", Sub: "sub", Aud: "aud", Scp: "scp", Iss: "iss", Iat: 1000, Exp: 999}
	assert.EqualError(t, ValidatePayloadFields(p), "invalid token: exp (999) must be greater than iat (1000)")
}

func TestValidatePayloadFields_ZeroExp_Valid(t *testing.T) {
	p := Payload{Jti: "jti", Sub: "sub", Aud: "aud", Scp: "scp", Iss: "iss", Iat: 1000}
	assert.NoError(t, ValidatePayloadFields(p))
}

func TestValidateExpiration_NotExpired(t *testing.T) {
	p := Payload{Exp: float64(time.Now().Add(1 * time.Hour).Unix())}
	assert.NoError(t, ValidateExpiration(p, time.Now()))
}

func TestValidateExpiration_Expired(t *testing.T) {
	p := Payload{Exp: float64(time.Now().Add(-1 * time.Hour).Unix())}
	assert.ErrorContains(t, ValidateExpiration(p, time.Now()), "token expired at")
}

func TestValidateExpiration_ZeroExp_NoExpiry(t *testing.T) {
	assert.NoError(t, ValidateExpiration(Payload{}, time.Now()))
}

func TestValidateExpiration_UsesNow(t *testing.T) {
	p := Payload{Exp: 1000}
	assert.NoError(t, ValidateExpiration(p, time.Unix(999, 0)))
	assert.EqualError(t, ValidateExpiration(p, time.Unix(1000, 0)), "token expired at 1970-01-01T00:16:40Z")
}

func TestParseRootCertificate_Valid(t *testing.T) {
	certB64, _ := generateTestCert(t)
	key, err := ParseRootCertificate(certB64)
	require.NoError(t, err)
	assert.NotNil(t, key)
}

func TestParseRootCertificate_InvalidBase64(t *testing.T) {
	_, err := ParseRootCertificate("!!!invalid!!!")
	assert.Error(t, err)
}

func TestParseRootCertificate_InvalidDER(t *testing.T) {
	_, err := ParseRootCertificate(base64.StdEncoding.EncodeToString([]byte("not a certificate")))
	assert.Error(t, err)
}

func TestParseRootCertificate_NonRSAKey(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	certDER, err := x509.CreateCertificate(rand.Reader, certTemplate(), certTemplate(), &ecKey.PublicKey, ecKey)
	require.NoError(t, err)
	_, err = ParseRootCertificate(base64.StdEncoding.EncodeToString(certDER))
	assert.ErrorContains(t, err, "not RSA")
}

func TestParseRootCertificate_TrimsWhitespace(t *testing.T) {
	certB64, _ := generateTestCert(t)
	_, err := ParseRootCertificate(certB64 + "\n")
	assert.NoError(t, err)
}

func TestPayload_Scopes(t *testing.T) {
	p := Payload{Scp: "applied-permissions/user system:metrics:r"}
	assert.Equal(t, []string{"applied-permissions/user", "system:metrics:r"}, p.Scopes())
}

func TestPayload_ScopesEmpty(t *testing.T) {
	assert.Nil(t, (&Payload{}).Scopes())
}

func TestPayload_AudAsStringSlice_String(t *testing.T) {
	assert.Equal(t, []string{"*@*"}, (&Payload{Aud: "*@*"}).AudAsStringSlice())
}

func TestPayload_AudAsStringSlice_Array(t *testing.T) {
	p := Payload{Aud: []interface{}{"svc1@id1", "svc2@id2"}}
	assert.Equal(t, []string{"svc1@id1", "svc2@id2"}, p.AudAsStringSlice())
}

func TestPayload_IatExpTime(t *testing.T) {
	p := Payload{Iat: 1677585465, Exp: 1677585465}
	assert.Equal(t, int64(1677585465), p.IatTime().Unix())
	assert.Equal(t, int64(1677585465), p.ExpTime().Unix())
}
