package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"strings"
	"time"
)

func extractPayloadFromAccessToken(token string) (TokenPayload, error) {
	// Separate token parts.
	tokenParts := strings.Split(token, ".")

	// Decode the payload.
	if len(tokenParts) != 3 {
		return TokenPayload{}, errorutils.CheckError(errors.New("received invalid access-token"))
	}
	payload, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return TokenPayload{}, errorutils.CheckError(err)
	}

	// Unmarshal json.
	var tokenPayload TokenPayload
	err = json.Unmarshal(payload, &tokenPayload)
	if err != nil {
		return TokenPayload{}, errorutils.CheckError(errors.New("Failed extracting payload from the provided access-token." + err.Error()))
	}
	err = setAudienceManually(&tokenPayload, payload)
	return tokenPayload, err
}

// Audience field was changed from string to string[]. This function extracts this field manually to allow backward compatibility.
func setAudienceManually(tokenPayload *TokenPayload, payload []byte) error {
	allValuesMap := make(map[string]interface{})
	err := json.Unmarshal(payload, &allValuesMap)
	if err != nil {
		return errorutils.CheckError(errors.New("Failed extracting audience from payload. " + err.Error()))
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
	return errorutils.CheckError(errors.New("failed extracting audience from payload. Audience is of unexpected type"))
}

func ExtractUsernameFromAccessToken(token string) (string, error) {
	tokenPayload, err := extractPayloadFromAccessToken(token)
	if err != nil {
		return "", err
	}
	// Extract subject.
	if tokenPayload.Subject == "" {
		return "", errorutils.CheckError(errors.New("could not extract subject from the provided access-token"))
	}

	// Extract username from subject.
	usernameStartIndex := strings.LastIndex(tokenPayload.Subject, "/")
	if usernameStartIndex < 0 {
		return "", errorutils.CheckError(errors.New(fmt.Sprintf("Could not extract username from access-token's subject: %s", tokenPayload.Subject)))
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
var RefreshBeforeExpiryMinutes = int64(10)

const WaitBeforeRefreshSeconds = 15
