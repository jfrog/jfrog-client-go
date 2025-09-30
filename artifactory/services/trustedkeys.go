package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// TrustedKeysService handles trusted keys operations
type TrustedKeysService struct {
	client         *jfroghttpclient.JfrogHttpClient
	serviceDetails auth.ServiceDetails
}

// TrustedKeyParams represents the parameters for trusted key operations
type TrustedKeyParams struct {
	Alias     string `json:"alias"`
	PublicKey string `json:"key"`
}

// TrustedKeyResponse represents the response from trusted keys API
type TrustedKeyResponse struct {
	Alias   string `json:"alias,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// TrustedKeyInfo represents a trusted key entry returned by GET
type TrustedKeyInfo struct {
	Kid         string `json:"kid"`
	Type        string `json:"type"`
	Alias       string `json:"alias"`
	Fingerprint string `json:"fingerprint"`
	IssuedBy    string `json:"issuedBy"`
	Issued      int64  `json:"issued"`
	Expiry      int64  `json:"expiry"`
}

// TrustedKeysResponse represents the response wrapper from GET /api/security/keys/trusted
type TrustedKeysResponse struct {
	Keys []TrustedKeyInfo `json:"keys"`
}

// NewTrustedKeysService creates a new TrustedKeysService instance
func NewTrustedKeysService(client *jfroghttpclient.JfrogHttpClient) *TrustedKeysService {
	return &TrustedKeysService{client: client}
}

// GetJfrogHttpClient returns the http client
func (tks *TrustedKeysService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return tks.client
}

// SetServiceDetails sets the service details
func (tks *TrustedKeysService) SetServiceDetails(serviceDetails auth.ServiceDetails) {
	tks.serviceDetails = serviceDetails
}

// UploadTrustedKey uploads a public key to the JFrog platform trusted keys
func (tks *TrustedKeysService) UploadTrustedKey(params TrustedKeyParams) (*TrustedKeyResponse, error) {
	if params.Alias == "" {
		return nil, errorutils.CheckErrorf("key alias cannot be empty")
	}
	if params.PublicKey == "" {
		return nil, errorutils.CheckErrorf("public key cannot be empty")
	}

	// Build the API URL
	requestUrl, err := tks.buildTrustedKeysUrl()
	if err != nil {
		return nil, err
	}

	// Prepare request body
	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	// Prepare HTTP client details
	httpClientsDetails := tks.serviceDetails.CreateHttpClientDetails()
	httpClientsDetails.Headers["Content-Type"] = "application/json"

	log.Info(fmt.Sprintf("Uploading trusted key with alias '%s' to JFrog platform...", params.Alias))

	// Send POST request
	resp, body, err := tks.client.SendPost(requestUrl, requestBody, &httpClientsDetails)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	// Parse response
	var response TrustedKeyResponse
	if len(body) > 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			log.Debug("Failed to parse response JSON, treating as raw message:", string(body))
			response.Message = string(body)
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errorMsg := fmt.Sprintf("trusted keys API returned status %d", resp.StatusCode)
		
		// Add specific error messages for common status codes
		switch resp.StatusCode {
		case http.StatusForbidden:
			errorMsg += ": Forbidden - insufficient permissions to upload trusted keys"
		case http.StatusNotFound:
			errorMsg += ": 404 page not found - trusted keys API endpoint not available"
		case http.StatusUnauthorized:
			errorMsg += ": Unauthorized - invalid or expired authentication token"
		}
		
		// Add response details to error message in priority order
		var additionalInfo string
		switch {
		case response.Error != "":
			additionalInfo = response.Error
		case response.Message != "":
			additionalInfo = response.Message
		case len(body) > 0:
			additionalInfo = string(body)
		}
		
		if additionalInfo != "" {
			errorMsg += ": " + additionalInfo
		}
		return &response, errorutils.CheckErrorf("%s", errorMsg)
	}

	log.Info(fmt.Sprintf("✓ Trusted key '%s' uploaded successfully to JFrog platform", params.Alias))
	return &response, nil
}

// buildTrustedKeysUrl builds the trusted keys API URL
func (tks *TrustedKeysService) buildTrustedKeysUrl() (string, error) {
	baseUrl := tks.serviceDetails.GetUrl()
	if baseUrl == "" {
		return "", errorutils.CheckErrorf("service URL cannot be empty")
	}

	// Build URL without query parameters
	requestUrl := utils.AddTrailingSlashIfNeeded(baseUrl) + "api/security/keys/trusted"
	return requestUrl, nil
}
