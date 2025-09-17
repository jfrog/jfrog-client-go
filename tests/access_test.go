//go:build itest

package tests

import (
	"testing"
)

func initAccessTest(t *testing.T) {
	if !*TestAccess {
		t.Skip("Skipping access test. To run access test add the '-test.access=true' option.")
	}
}
