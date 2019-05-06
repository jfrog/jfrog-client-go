package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"strings"
)

func ExtractUsernameFromAccessToken(token string) (string, error) {
	// Separate token parts.
	tokenParts := strings.Split(token, ".")

	// Decode the payload.
	if len(tokenParts) != 3 {
		return "", errorutils.CheckError(errors.New("Received invalid access-token."))
	}
	payload, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return "", errorutils.CheckError(err)
	}

	// Unmarshal json.
	var tokenPayload tokenPayload
	err = json.Unmarshal(payload, &tokenPayload)
	if err != nil {
		return "", errorutils.CheckError(errors.New("Failed extracting payload from the provided access-token." + err.Error()))
	}

	// Extract subject.
	if tokenPayload.Subject == "" {
		return "", errorutils.CheckError(errors.New("Could not extract subject from the provided access-token."))
	}

	// Extract username from subject.
	usernameStartIndex := strings.LastIndex(tokenPayload.Subject, "/")
	if usernameStartIndex < 0 {
		return "", errorutils.CheckError(errors.New(fmt.Sprintf("Could not extract username from access-token's subject: %s", tokenPayload.Subject)))
	}
	username := tokenPayload.Subject[usernameStartIndex+1:]

	return username, nil
}

type tokenPayload struct {
	Subject        string `json:"sub,omitempty"`
	Scope          string `json:"scp,omitempty"`
	Audience       string `json:"aud,omitempty"`
	Issuer         string `json:"iss,omitempty"`
	ExpirationTime int    `json:"exp,omitempty"`
	IssuedAt       int    `json:"iat,omitempty"`
	JwtId          string `json:"jti,omitempty"`
}
