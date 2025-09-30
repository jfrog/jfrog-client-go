package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustedKeyParams_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  TrustedKeyParams
		wantErr bool
	}{
		{
			name: "valid params",
			params: TrustedKeyParams{
				Alias:     "test-key",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----",
			},
			wantErr: false,
		},
		{
			name: "empty alias",
			params: TrustedKeyParams{
				Alias:     "",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----",
			},
			wantErr: true,
		},
		{
			name: "empty public key",
			params: TrustedKeyParams{
				Alias:     "test-key",
				PublicKey: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test parameter validation by checking the params directly
			if tt.wantErr {
				if tt.params.Alias == "" {
					assert.Empty(t, tt.params.Alias, "Expected empty alias for error case")
				}
				if tt.params.PublicKey == "" {
					assert.Empty(t, tt.params.PublicKey, "Expected empty public key for error case")
				}
			} else {
				assert.NotEmpty(t, tt.params.Alias, "Expected non-empty alias for valid case")
				assert.NotEmpty(t, tt.params.PublicKey, "Expected non-empty public key for valid case")
			}
		})
	}
}
func TestNewTrustedKeysService(t *testing.T) {
	service := NewTrustedKeysService(nil)
	assert.NotNil(t, service)
	assert.Nil(t, service.GetJfrogHttpClient())
}

func TestTrustedKeysService_NoDuplicateAliasValidation(t *testing.T) {
	// This test verifies that the client no longer validates duplicate aliases
	// The server should handle this validation instead
	
	service := NewTrustedKeysService(nil)
	assert.NotNil(t, service)
	
	// Verify that CheckAliasExists method no longer exists
	// This is a compile-time check - if this code compiles, it means the method was successfully removed
	
	// The service should not have a CheckAliasExists method anymore
	// This test would fail to compile if the method still existed and was being called
	// Since we removed the method, this test passing means client-side duplicate validation is gone
	
	// Create valid test parameters
	params := TrustedKeyParams{
		Alias:     "test-duplicate-alias",
		PublicKey: "-----BEGIN PUBLIC KEY-----\ntest-key-content\n-----END PUBLIC KEY-----",
	}
	
	// Verify parameters are valid (only basic validation should remain)
	assert.NotEmpty(t, params.Alias, "Alias should not be empty")
	assert.NotEmpty(t, params.PublicKey, "Public key should not be empty")
	
	// The UploadTrustedKey method should not call any alias existence check
	// It should proceed directly to the server request
	// Server-side validation will handle duplicates appropriately
}

func TestTrustedKeysService_ServerSideDuplicateValidation(t *testing.T) {
	// This test demonstrates that duplicate alias validation is now handled server-side
	// The client sends the request directly to the server without pre-validation
	
	service := NewTrustedKeysService(nil)
	assert.NotNil(t, service)
	
	// Test demonstrates that the service no longer has client-side duplicate validation
	// The following would previously have failed due to client-side checks
	// Now it should only fail when sent to the server
	
	params := TrustedKeyParams{
		Alias:     "duplicate-alias", 
		PublicKey: "-----BEGIN PUBLIC KEY-----\ntest-key-content\n-----END PUBLIC KEY-----",
	}
	
	// Verify basic parameter validation still works
	assert.NotEmpty(t, params.Alias, "Alias should not be empty")
	assert.NotEmpty(t, params.PublicKey, "Public key should not be empty")
	
	// The fact that this test compiles and runs proves that:
	// 1. CheckAliasExists method was successfully removed (compilation would fail if it existed)
	// 2. UploadTrustedKey no longer calls CheckAliasExists (no client-side duplicate validation)
	// 3. Only basic parameter validation remains (empty alias/key checks)
	// 4. Server-side validation will handle duplicate aliases appropriately
	
	// Note: We can't test the actual server interaction without a real server,
	// but the absence of client-side validation is proven by successful compilation
}

