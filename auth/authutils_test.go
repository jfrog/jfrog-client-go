package auth

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
)

var (
	// #nosec G101 -- Dummy tokens for tests jfrog-ignore
	token1 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0.eyJzdWIiOiJqZnJ0QDAxYzNnZmZoZzJlOHc2MTQ5ZTNhMnEwdzk3XC91c2Vyc1wvYWRtaW4iLCJzY3AiOiJtZW1iZXItb2YtZ3JvdXBzOnJlYWRlcnMgYXBpOioiLCJhdWQiOiJqZnJ0QDAxYzNnZmZoZzJlOHc2MTQ5ZTNhMnEwdzk3IiwiaXNzIjoiamZydEAwMWMzZ2ZmaGcyZTh3NjE0OWUzYTJxMHc5NyIsImV4cCI6MTU1NjAzNzc2NSwiaWF0IjoxNTU2MDM0MTY1LCJqdGkiOiI1M2FlMzgyMy05NGM3LTQ0OGItOGExOC1iZGVhNDBiZjFlMjAifQ.Bp3sdvppvRxysMlLgqT48nRIHXISj9sJUCXrm7pp8evJGZW1S9hFuK1olPmcSybk2HNzdzoMcwhUmdUzAssiQkQvqd_HanRcfFbrHeg5l1fUQ397ECES-r5xK18SYtG1VR7LNTVzhJqkmRd3jzqfmIK2hKWpEgPfm8DRz3j4GGtDRxhb3oaVsT2tSSi_VfT3Ry74tzmO0GcCvmBE2oh58kUZ4QfEsalgZ8IpYHTxovsgDx_M7ujOSZx_hzpz-iy268-OkrU22PQPCfBmlbEKeEUStUO9n0pj4l1ODL31AGARyJRy46w4yzhw7Fk5P336WmDMXYs5LAX2XxPFNLvNzA"
	// #nosec G101 jfrog-ignore
	token2 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0.eyJzdWIiOiJqZnJ0QDAwMWMzZ2ZmaGcyZTh3NjE0OWUzYTJxMHc5NyIsImV4cCI6MTU1NjAzNzc2NSwiaWF0IjoxNTU2MDM0MTY1LCJqdGkiOiI1M2FlMzgyMy05NGM3LTQ0OGItOGExOC1iZGVhNDBiZjFlMjAifQ.Bp3sdvppvRxysMlLgqT48nRIHXISj9sJUCXrm7pp8evJGZW1S9hFuK1olPmcSybk2HNzdzoMcwhUmdUzAssiQkQvqd_HanRcfFbrHeg5l1fUQ397ECES-r5xK18SYtG1VR7LNTVzhJqkmRd3jzqfmIK2hKWpEgPfm8DRz3j4GGtDRxhb3oaVsT2tSSi_VfT3Ry74tzmO0GcCvmBE2oh58kUZ4QfEsalgZ8IpYHTxovsgDx_M7ujOSZx_hzpz-iy268-OkrU22PQPCfBmlbEKeEUStUO9n0pj4l1ODL31AGARyJRy46w4yzhw7Fk5P336WmDMXYs5LAX2XxPFNLvNzA"
	// #nosec G101 jfrog-ignore
	token3 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0"
	// #nosec G101 jfrog-ignore
	token4 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJsS0NYXzFvaTBQbTZGdF9XRklkejZLZ1g4U0FULUdOY0lJWXRjTC1KM084In0.eyJzdWIiOiJqZmZlQDAwMFwvdXNlcnNcL3RlbXB1c2VyIiwic2NwIjoiYXBwbGllZC1wZXJtaXNzaW9uc1wvYWRtaW4gYXBpOioiLCJhdWQiOlsiamZydEAqIiwiamZtZEAqIiwiamZldnRAKiIsImpmYWNAKiJdLCJpc3MiOiJqZmZlQDAwMCIsImV4cCI6MTYxNjQ4OTU4NSwiaWF0IjoxNjE2NDg1OTg1LCJqdGkiOiI0OTBlYWEzOS1mMzYxLTQxYjAtOTA5Ni1kNjg5NmQ0ZWQ3YjEifQ.J5P8Pu5tqEjgnLFLEoCdh1LJHWiMmEHht95v0EFuixwO-osq7sfXua_UCGBkKbmqVSGKew9kl_LTcbq_uMe281_5q2yYxT74iqc2wQ1K0uovEUeIU6E65oi70JwUWUwcF3sNJ2gFatnvgSu-2Kv6m-DtSIW36WS3Mh8uMZQ19ob4fmueVmMFyQsp0EEG6xFYeOK6SB8OUd0gAd_XvXiSRuF0eLabhKmXM2pVBLYfd2KIMlkFckEOGGOzeglvA62xmP4Ik7UsF487NAo0LeS_Pd79owr0jtgTYkCTrLlFhUzUMDVmD_LsCMyf_S4CJxhwkCRhhy9SYSs1WPgknL3--w"
	// #nosec G101 jfrog-ignore
	token5 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJDRlVIRER4UXZaM1VNWEZxS0ZWUlFiOEFROEd6VWxSZkZJMEt4TmIzdk1jIn0.eyJzdWIiOiJ5YWhhdi90ZXN0LXJlcG8iLCJzY3AiOiJhcHBsaWVkLXBlcm1pc3Npb25zL2dyb3VwczpcImFkbWluLWdyb3VwXCIsIiwiYXVkIjoiKkAqIiwiaXNzIjoiamZhY0AwMWdnbXFxcDc0MzZuOTB3d3I4Ym5nMXp5OSIsImV4cCI6MTcxNTE4MzA3MiwiaWF0IjoxNzE1MTgzMDEyLCJqdGkiOiJmN2IxMmIzMi0xMmNkLTQ1Y2ItYWZjYS1iNTYyMjc2YjY0YmQifQ.I6df8E0_1t7uYzSQkiQBNh9GIGyr541rIRQ8BDD401N4DV98dWsqACmdlYTOAaxn_el4Lz7_OaK0GnVNGf9hiZz9QKaXbe-HnL9jG-TobpOlyhkc6iBpnizuZ9T9YiveCG_NgDMWn5syiZ912t6PuZqNN2JmwswqfE9QDm6xCH8fu7h0Rs1qDNkahtgQvO99e5d7LnuOS9VfkXBxLDZ5AeUbd89zmujgDB4hMXB3J-dQ3QxGHRPS_KUo7sf7lRvn4PydYkhbhrhg6GP6ss6rMmEJM5v8azMTrkLwksoCWtK9YpD5S70f7AoE5U5j5BttZ0S5dPGagKWZJiA1egna-w"
)

func TestExtractUsernameFromAccessToken(t *testing.T) {
	testCases := []struct {
		expectedUsername string
		inputToken       string
		shouldError      bool
	}{
		{"admin", token1, false},
		{"", token2, true},
		{"", token3, true},
		{"tempuser", token4, false},
		{"yahav/test-repo", token5, false},
	}
	// Discard output logging to prevent negative logs
	previousLog := tests.RedirectLogOutputToNil()
	defer func() {
		log.SetLogger(previousLog)
	}()

	for _, testCase := range testCases {
		username := ExtractUsernameFromAccessToken(testCase.inputToken)
		assert.Equal(t, testCase.shouldError, username == "")
		assert.Equal(t, testCase.expectedUsername, username)
	}
}

func TestExtractSubjectFromAccessToken(t *testing.T) {
	testCases := []struct {
		expectedSubject string
		inputToken      string
		shouldError     bool
	}{
		{"jfrt@01c3gffhg2e8w6149e3a2q0w97/users/admin", token1, false},
		{"jfrt@001c3gffhg2e8w6149e3a2q0w97", token2, false},
		{"", token3, true},
		{"jffe@000/users/tempuser", token4, false},
		{"yahav/test-repo", token5, false},
	}

	for _, testCase := range testCases {
		subject, err := ExtractSubjectFromAccessToken(testCase.inputToken)
		if testCase.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, testCase.expectedSubject, subject)
	}
}
