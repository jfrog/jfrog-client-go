package tests

import (
	"fmt"
	accessservices "github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessInvite(t *testing.T) {
	initAccessTest(t)
	t.Run("invite", testInviteUser)
}

func testInviteUser(t *testing.T) {
	randomMail := fmt.Sprintf("test%s@jfrog.com", timestampStr)
	UserParams := getTestInvitedUserParams(randomMail)
	err := testUserService.CreateUser(UserParams)
	assert.NoError(t, err)

	// Second invitation should fail because we can invite user only once a day (Access's internal reasons).
	err = testsAccessInviteService.InviteUser(randomMail, "cli")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "already invited today"), "error : "+err.Error())
}

func getTestInvitedUserParams(email string) services.UserParams {
	// Data members "name" and "email" should both be the email (Access's internal reasons).
	userDetails := services.User{
		Name:                     email,
		Email:                    email,
		Password:                 "Password1!",
		Admin:                    &trueValue,
		ShouldInvite:             &trueValue,
		Source:                   accessservices.InviteCliSourceName,
		ProfileUpdatable:         &trueValue,
		DisableUIAccess:          &falseValue,
		InternalPasswordDisabled: &falseValue,
		WatchManager:             &trueValue,
		ReportsManager:           &trueValue,
		PolicyManager:            &trueValue,
		ProjectAdmin:             &trueValue,
	}
	return services.UserParams{
		UserDetails:     userDetails,
		ReplaceIfExists: false,
	}
}
