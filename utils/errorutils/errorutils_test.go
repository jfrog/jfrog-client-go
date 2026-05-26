package errorutils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckResponseStatusWithBody_ReturnsTypedHttpResponseError(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusUnauthorized, Status: "401 Unauthorized"}
	body := []byte(`{"errors":[{"code":"UNAUTHORIZED","message":"bad credentials"}]}`)

	err := CheckResponseStatusWithBody(resp, body, http.StatusOK)
	assert.Error(t, err)

	var httpErr *HttpResponseError
	assert.True(t, errors.As(err, &httpErr), "expected error to unwrap to *HttpResponseError")
	assert.Equal(t, http.StatusUnauthorized, httpErr.StatusCode)
	assert.Equal(t, "401 Unauthorized", httpErr.Status)
	assert.Equal(t, body, httpErr.Body)

	// Error() preserves the legacy human-readable format.
	assert.True(t, strings.HasPrefix(err.Error(), "server response: 401 Unauthorized\n"),
		"legacy text format must be preserved, got: %q", err.Error())
}

func TestCheckResponseStatusWithBody_HappyPathReturnsNil(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusOK, Status: "200 OK"}
	assert.NoError(t, CheckResponseStatusWithBody(resp, []byte(`{}`), http.StatusOK))
}

func TestCheckResponseStatus_ReadsBodyAndReturnsTyped(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       io.NopCloser(strings.NewReader("not-json")),
	}
	err := CheckResponseStatus(resp, http.StatusOK)
	assert.Error(t, err)

	var httpErr *HttpResponseError
	assert.True(t, errors.As(err, &httpErr))
	assert.Equal(t, 500, httpErr.StatusCode)
	assert.Equal(t, []byte("not-json"), httpErr.Body)
}

func TestHttpResponseError_UnwrapsThroughFmtErrorf(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusForbidden, Status: "403 Forbidden"}
	inner := CheckResponseStatusWithBody(resp, []byte(`{"msg":"nope"}`), http.StatusOK)
	wrapped := fmt.Errorf("failed to exchange OIDC token: %w", inner)

	var httpErr *HttpResponseError
	assert.True(t, errors.As(wrapped, &httpErr),
		"errors.As must walk through fmt.Errorf %%w wrapping (OIDC code path)")
	assert.Equal(t, 403, httpErr.StatusCode)
}
