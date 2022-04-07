package tests

import (
	"fmt"
	accessservices "github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessInvite(t *testing.T) {
	initAccessTest(t)
	t.Run("invite", testInviteUser)
}

func testInviteUser(t *testing.T) {
	randomMail := fmt.Sprintf("test%s@jfrog.com", timestampStr)
	randomMail = "gaimacrame@gmail.com"
	UserParams := getTestInvitedUserParams(randomMail)
	err := testUserService.CreateUser(UserParams)
	assert.NoError(t, err)
	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Second invitation should fail because we can invite user only once a day for access internal reasons.
	err = testsAccessInviteService.InviteUser(randomMail)
	assert.Error(t, err)
	// TODO: check error
	//assert.True(t, strings.HasSuffix(err.Error(), "already invited today"), "error : "+err.Error())
	//err = testUserService.DeleteUser(UserParams.UserDetails.Name)
	assert.NoError(t, err)
}

func getTestInvitedUserParams(email string) services.UserParams {
	// Data members "name" and "email" should both be the email for internal access reasons.
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
