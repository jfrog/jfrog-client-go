package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* #nosec G101 -- Dummy tokens for tests. */
var (
	token1 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0.eyJzdWIiOiJqZnJ0QDAxYzNnZmZoZzJlOHc2MTQ5ZTNhMnEwdzk3XC91c2Vyc1wvYWRtaW4iLCJzY3AiOiJtZW1iZXItb2YtZ3JvdXBzOnJlYWRlcnMgYXBpOioiLCJhdWQiOiJqZnJ0QDAxYzNnZmZoZzJlOHc2MTQ5ZTNhMnEwdzk3IiwiaXNzIjoiamZydEAwMWMzZ2ZmaGcyZTh3NjE0OWUzYTJxMHc5NyIsImV4cCI6MTU1NjAzNzc2NSwiaWF0IjoxNTU2MDM0MTY1LCJqdGkiOiI1M2FlMzgyMy05NGM3LTQ0OGItOGExOC1iZGVhNDBiZjFlMjAifQ.Bp3sdvppvRxysMlLgqT48nRIHXISj9sJUCXrm7pp8evJGZW1S9hFuK1olPmcSybk2HNzdzoMcwhUmdUzAssiQkQvqd_HanRcfFbrHeg5l1fUQ397ECES-r5xK18SYtG1VR7LNTVzhJqkmRd3jzqfmIK2hKWpEgPfm8DRz3j4GGtDRxhb3oaVsT2tSSi_VfT3Ry74tzmO0GcCvmBE2oh58kUZ4QfEsalgZ8IpYHTxovsgDx_M7ujOSZx_hzpz-iy268-OkrU22PQPCfBmlbEKeEUStUO9n0pj4l1ODL31AGARyJRy46w4yzhw7Fk5P336WmDMXYs5LAX2XxPFNLvNzA"
	token2 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0.eyJzdWIiOiJqZnJ0QDAwMWMzZ2ZmaGcyZTh3NjE0OWUzYTJxMHc5NyIsImV4cCI6MTU1NjAzNzc2NSwiaWF0IjoxNTU2MDM0MTY1LCJqdGkiOiI1M2FlMzgyMy05NGM3LTQ0OGItOGExOC1iZGVhNDBiZjFlMjAifQ.Bp3sdvppvRxysMlLgqT48nRIHXISj9sJUCXrm7pp8evJGZW1S9hFuK1olPmcSybk2HNzdzoMcwhUmdUzAssiQkQvqd_HanRcfFbrHeg5l1fUQ397ECES-r5xK18SYtG1VR7LNTVzhJqkmRd3jzqfmIK2hKWpEgPfm8DRz3j4GGtDRxhb3oaVsT2tSSi_VfT3Ry74tzmO0GcCvmBE2oh58kUZ4QfEsalgZ8IpYHTxovsgDx_M7ujOSZx_hzpz-iy268-OkrU22PQPCfBmlbEKeEUStUO9n0pj4l1ODL31AGARyJRy46w4yzhw7Fk5P336WmDMXYs5LAX2XxPFNLvNzA"
	token3 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJIcnU2VHctZk1yOTV3dy12TDNjV3ZBVjJ3Qm9FSHpHdGlwUEFwOE1JdDljIn0"
	token4 = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJsS0NYXzFvaTBQbTZGdF9XRklkejZLZ1g4U0FULUdOY0lJWXRjTC1KM084In0.eyJzdWIiOiJqZmZlQDAwMFwvdXNlcnNcL3RlbXB1c2VyIiwic2NwIjoiYXBwbGllZC1wZXJtaXNzaW9uc1wvYWRtaW4gYXBpOioiLCJhdWQiOlsiamZydEAqIiwiamZtZEAqIiwiamZldnRAKiIsImpmYWNAKiJdLCJpc3MiOiJqZmZlQDAwMCIsImV4cCI6MTYxNjQ4OTU4NSwiaWF0IjoxNjE2NDg1OTg1LCJqdGkiOiI0OTBlYWEzOS1mMzYxLTQxYjAtOTA5Ni1kNjg5NmQ0ZWQ3YjEifQ.J5P8Pu5tqEjgnLFLEoCdh1LJHWiMmEHht95v0EFuixwO-osq7sfXua_UCGBkKbmqVSGKew9kl_LTcbq_uMe281_5q2yYxT74iqc2wQ1K0uovEUeIU6E65oi70JwUWUwcF3sNJ2gFatnvgSu-2Kv6m-DtSIW36WS3Mh8uMZQ19ob4fmueVmMFyQsp0EEG6xFYeOK6SB8OUd0gAd_XvXiSRuF0eLabhKmXM2pVBLYfd2KIMlkFckEOGGOzeglvA62xmP4Ik7UsF487NAo0LeS_Pd79owr0jtgTYkCTrLlFhUzUMDVmD_LsCMyf_S4CJxhwkCRhhy9SYSs1WPgknL3--w"
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
	}

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
