package tests

import (
	"fmt"
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
	UserParams := getTestInvitedUserParams(false, randomMail)
	err := testUserService.CreateUser(UserParams)
	assert.NoError(t, err)

	// TDO: why create user?
	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)
	assert.Nil(t, user)
	err = testsAccessInviteService.InviteUser(randomMail)
	assert.NoError(t, err)
	err = testsAccessInviteService.InviteUser(randomMail)
	assert.Error(t, err)
	err = testUserService.DeleteUser(UserParams.UserDetails.Name)
	assert.NoError(t, err)
}

func getTestInvitedUserParams(replaceIfExists bool, email string) services.UserParams {
	userDetails := services.User{
		Name:                     email,
		Email:                    email,
		Password:                 "Password1!",
		Admin:                    &trueValue,
		Realm:                    "internal",
		ProfileUpdatable:         &trueValue,
		DisableUIAccess:          &falseValue,
		InternalPasswordDisabled: &falseValue,
		ShouldInvite:             &trueValue,
		Source:                   services.InviteCliSourceName,
	}
	return services.UserParams{
		UserDetails:     userDetails,
		ReplaceIfExists: replaceIfExists,
	}
}
