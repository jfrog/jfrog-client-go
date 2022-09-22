package auth

import (
	"encoding/base64"
	"encoding/json"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"strings"
	"time"
)

type CreateTokenResponseData struct {
	CommonTokenParams
}

type CommonTokenParams struct {
	Scope        string `json:"scope,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Refreshable  *bool  `json:"refreshable,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	Audience     string `json:"audience,omitempty"`
}

func extractPayloadFromAccessToken(token string) (TokenPayload, error) {
	// Separate token parts.
	tokenParts := strings.Split(token, ".")

	// Decode the payload.
	if len(tokenParts) != 3 {
		return TokenPayload{}, errorutils.CheckErrorf("couldn't extract payload from Access Token.\n" +
			"Hint: Reference and Identity Tokens are currently not supported by this functionality. " +
			"You can use them as the password in Basic Authentication (username and password). " +
			"(supported by JFrog Artifactory 7.43.0 or higher)",
		)
	}
	payload, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return TokenPayload{}, errorutils.CheckError(err)
	}

	// Unmarshal json.
	var tokenPayload TokenPayload
	err = json.Unmarshal(payload, &tokenPayload)
	if err != nil {
		return TokenPayload{}, errorutils.CheckErrorf("failed extracting payload from the provided access-token. " + err.Error())
	}
	err = setAudienceManually(&tokenPayload, payload)
	return tokenPayload, err
}

// Audience field was changed from string to string[]. This function extracts this field manually to allow backward compatibility.
func setAudienceManually(tokenPayload *TokenPayload, payload []byte) error {
	allValuesMap := make(map[string]interface{})
	err := json.Unmarshal(payload, &allValuesMap)
	if err != nil {
		return errorutils.CheckErrorf("Failed extracting audience from payload. " + err.Error())
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

func ExtractUsernameFromAccessToken(token string) (string, error) {
	tokenPayload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return "", err
	}
	// Extract subject.
	if tokenPayload.Subject == "" {
		return "", errorutils.CheckErrorf("could not extract subject from the provided access-token")
	}

	// Extract username from subject.
	usernameStartIndex := strings.LastIndex(tokenPayload.Subject, "/")
	if usernameStartIndex < 0 {
		return "", errorutils.CheckErrorf("Could not extract username from access-token's subject: %s", tokenPayload.Subject)
	}
	username := tokenPayload.Subject[usernameStartIndex+1:]

	return username, nil
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

// RefreshBeforeExpiryMinutes Artifactory's refresh token mechanism creates tokens that expired in 60 minutes. We want to refresh them after 50 minutes (when 10 minutes left)
var RefreshBeforeExpiryMinutes = int64(10)

// InviteRefreshBeforeExpiryMinutes Invitations mechanism creates tokens that are valid for 1 year. We want to refresh the token every 50 minutes.
var InviteRefreshBeforeExpiryMinutes = int64(365*24*60 - 50)

const WaitBeforeRefreshSeconds = 15
