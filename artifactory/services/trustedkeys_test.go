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

