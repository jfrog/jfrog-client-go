package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/http/httpclient"

	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type CreateTokenResponseData struct {
	CommonTokenParams
	ReferenceToken string `json:"reference_token,omitempty"`
	TokenId        string `json:"token_id,omitempty"`
}

type OidcTokenResponseData struct {
	CommonTokenParams
	IssuedTokenType string `json:"issued_token_type,omitempty"`
	Username        string `json:"username,omitempty"`
}

type CommonTokenParams struct {
	Scope        string `json:"scope,omitempty"`
	AccessToken  string `json:"access_token,omitempty"` // #nosec G117 -- API struct for OAuth token response
	ExpiresIn    *uint  `json:"expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Refreshable  *bool  `json:"refreshable,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"` // #nosec G117 -- API struct for OAuth token response
	GrantType    string `json:"grant_type,omitempty"`
	Audience     string `json:"audience,omitempty"`
}

func extractPayloadFromAccessToken(token string) (TokenPayload, error) {
	// Separate token parts.
	tokenParts := strings.Split(token, ".")

	// Decode the payload.
	if len(tokenParts) != 3 {
		return TokenPayload{}, errorutils.CheckErrorf("couldn't extract payload from Access Token")
	}
	payload, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return TokenPayload{}, errorutils.CheckError(err)
	}

	// Unmarshal json.
	var tokenPayload TokenPayload
	err = json.Unmarshal(payload, &tokenPayload)
	if err != nil {
		return TokenPayload{}, errorutils.CheckErrorf("failed extracting payload from the provided access-token: %s", err.Error())
	}
	err = setAudienceManually(&tokenPayload, payload)
	return tokenPayload, err
}

// Audience field was changed from string to string[]. This function extracts this field manually to allow backward compatibility.
func setAudienceManually(tokenPayload *TokenPayload, payload []byte) error {
	allValuesMap := make(map[string]interface{})
	err := json.Unmarshal(payload, &allValuesMap)
	if err != nil {
		return errorutils.CheckErrorf("failed extracting audience from payload: %s", err.Error())
	}
	aud, exists := allValuesMap["aud"]
	if !exists {
		return nil
	}
	if audStr, ok := aud.(string); ok {
		tokenPayload.Audience = audStr
		return nil
	}
	if audArray, ok := aud.([]interface{}); ok {
		for _, v := range audArray {
			if newAud, ok := v.(string); ok {
				tokenPayload.AudienceArray = append(tokenPayload.AudienceArray, newAud)
			}
		}
		return nil
	}
	return errorutils.CheckErrorf("failed extracting audience from payload. Audience is of unexpected type")
}

func ExtractUsernameFromAccessToken(token string) (username string) {
	if httpclient.IsApiKey(token) {
		log.Warn("The provided access token is an API key which should be used as password.")
		return
	}
	var err error
	defer func() {
		if err != nil {
			log.Warn(err.Error() + ".\n" +
				"The provided access token is not a valid JWT, probably a reference token.\n" +
				"Some package managers only support basic authentication which requires also a username.\n" +
				"If you plan to work with one of those package managers, please provide a username.")
		}
	}()
	tokenPayload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return
	}
	// Extract subject.
	if tokenPayload.Subject == "" {
		err = errorutils.CheckErrorf("couldn't extract subject from the provided access-token")
		return
	}

	// Extract username from subject.
	if strings.HasPrefix(tokenPayload.Subject, "jfrt@") || strings.Contains(tokenPayload.Subject, "/users/") {
		usernameStartIndex := strings.LastIndex(tokenPayload.Subject, "/")
		if usernameStartIndex < 0 {
			err = errorutils.CheckErrorf("couldn't extract username from access-token's subject: %s", tokenPayload.Subject)
			return
		}
		username = tokenPayload.Subject[usernameStartIndex+1:]
	} else {
		// OICD token for groups scope
		username = tokenPayload.Subject
	}
	if username == "" {
		err = errorutils.CheckErrorf("empty username extracted from access-token's subject: %s", tokenPayload.Subject)
	}
	return
}

// Extracts the expiry from an access token, in seconds
func ExtractExpiryFromAccessToken(token string) (int, error) {
	tokenPayload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return -1, err
	}
	expiry := tokenPayload.ExpirationTime - tokenPayload.IssuedAt
	return expiry, nil
}

// Returns 0 if expired
func GetTokenMinutesLeft(token string) (int64, error) {
	payload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return -1, err
	}
	left := int64(payload.ExpirationTime) - time.Now().Unix()
	if left < 0 {
		return 0, nil
	}
	return left / 60, nil
}

// Extracts the subject from an access token
func ExtractSubjectFromAccessToken(token string) (string, error) {
	tokenPayload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return "", err
	}
	return tokenPayload.Subject, nil
}

type TokenPayload struct {
	Subject        string `json:"sub,omitempty"`
	Scope          string `json:"scp,omitempty"`
	Issuer         string `json:"iss,omitempty"`
	ExpirationTime int    `json:"exp,omitempty"`
	IssuedAt       int    `json:"iat,omitempty"`
	JwtId          string `json:"jti,omitempty"`
	// Audience was changed to slice. Handle this field manually so the unmarshalling will not fail.
	Audience      string
	AudienceArray []string
}

// Refreshable Tokens Constants.
// Artifactory's refresh token mechanism creates tokens that expire in 60 minutes. We want to refresh them when 10 minutes are left.
var RefreshArtifactoryTokenBeforeExpiryMinutes = int64(10)

// Platform's access tokens are historically created with 1 year expiry, refreshed a week before expiry.
// Kept as the upper cap for CalculateProportionalRefreshThreshold, so tokens with long lifetimes keep this behavior.
var RefreshPlatformTokenBeforeExpiryMinutes = int64(7 * 24 * 60)

// Never refresh a platform token more than this many minutes before it expires, regardless of its lifetime.
const MinRefreshBeforeExpiryMinutes = int64(30)

// Refresh a platform token once this fraction of its total lifetime remains.
const refreshThresholdFraction = 0.10

const WaitBeforeRefreshSeconds = 15

// CalculateProportionalRefreshThreshold returns the number of minutes before expiry that a platform
// access token should be refreshed, proportional to its total lifetime. This avoids refreshing on
// every command for tokens issued with a shorter-than-1-year lifetime (e.g. 7 days), since the fixed
// RefreshPlatformTokenBeforeExpiryMinutes threshold would otherwise exceed the token's entire lifetime.
func CalculateProportionalRefreshThreshold(tokenLifetimeMinutes int64) int64 {
	threshold := int64(float64(tokenLifetimeMinutes) * refreshThresholdFraction)
	if threshold < MinRefreshBeforeExpiryMinutes {
		return MinRefreshBeforeExpiryMinutes
	}
	if threshold > RefreshPlatformTokenBeforeExpiryMinutes {
		return RefreshPlatformTokenBeforeExpiryMinutes
	}
	return threshold
}

// GetPlatformTokenRefreshThreshold extracts a platform access token's total lifetime (exp - iat) and
// returns its proportional refresh threshold, in minutes. See CalculateProportionalRefreshThreshold.
func GetPlatformTokenRefreshThreshold(token string) (int64, error) {
	expirySeconds, err := ExtractExpiryFromAccessToken(token)
	if err != nil {
		return 0, err
	}
	return CalculateProportionalRefreshThreshold(int64(expirySeconds) / 60), nil
}
