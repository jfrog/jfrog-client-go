package tests

import (
	clientUtils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/xsc/services/utils"
	"testing"
)

func TestXscVersion(t *testing.T) { // TODO eran this is skipped
	initXscTest(t, "", "")
	version, err := GetXscDetailsBackwardsCompatible(t).GetVersion()
	if err != nil {
		t.Error(err)
	}
	if version == "" {
		t.Error("Expected a version, got empty string")
	}
}

func initXscTest(t *testing.T, minXscVersion string, minXrayVersion string) {
	if !*TestXsc {
		t.Skip("Skipping xsc test. To run xsc test add the '-test.xsc=true' option.")
	}
	validateXscAndXrayVersion(t, minXscVersion, minXrayVersion)
}

// This func validates minimal Xsc version.
// Since Xsc was migrated into Xray from version 3.107.13, we need to check minimal Xray version from this version forward instead of Xsc version.
// For features that are available before the migration we pass minXscVersion to check. If the utilized Xray version >= 3.107.13, the returned Xsc version will always suffice the check.
// For features that were introduced only after the migration we pass only minXrayVersion to check and can leave minXscVersion blank.
// In general minXscVersion should be provided only for features that were introduced before Xsc migration to Xray
func validateXscAndXrayVersion(t *testing.T, minXscVersion string, minXrayVersion string) {
	// We first validate our Xray version so we will not address the old Xsc endpoints if Xray version >= 3.107.13. This will lead to a failure and skip the test
	currentXrayVersion, err := GetXrayDetails().GetVersion()
	if err != nil {
		t.Skip(err)
	}

	if minXrayVersion != "" {
		if err = clientUtils.ValidateMinimumVersion(clientUtils.Xsc, currentXrayVersion, minXrayVersion); err != nil {
			t.Skip(err)
		}
	}

	afterMigration := true
	if err = clientUtils.ValidateMinimumVersion(clientUtils.Xray, currentXrayVersion, utils.MinXrayVersionXscTransitionToXray); err != nil {
		err = nil
		afterMigration = false
	}
	if !afterMigration {
		// If Xray version < 3.107.13 we validate active Xsc server with minimal required version
		var currentXscVersion string
		currentXscVersion, err = GetXscDetails().GetVersion()
		if err != nil {
			t.Skip(err)
		}

		if minXscVersion != "" {
			if err = clientUtils.ValidateMinimumVersion(clientUtils.Xsc, currentXscVersion, minXscVersion); err != nil {
				t.Skip(err)
			}
		}
	}
}
