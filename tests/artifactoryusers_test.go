package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	initArtifactoryTest(t)
	t.Run("create", testCreateUser)
	t.Run("update", testUpdateUser)
	t.Run("clear users groups", testClearUserGroups)
	t.Run("delete", testDeleteUser)
}

func testCreateUser(t *testing.T) {
	UserParams := getTestUserParams(false, "")

	err := testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)
	// we don't know the default group when created, so just set it
	UserParams.UserDetails.Groups = user.Groups
	// password is not carried in reply
	user.Password = UserParams.UserDetails.Password
	assert.NotNil(t, user)
	assert.True(t, reflect.DeepEqual(UserParams.UserDetails, *user))
}

func testUpdateUser(t *testing.T) {
	UserParams := getTestUserParams(true, "")

	err := testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	UserParams.UserDetails.Email = "changed@mail.com"
	UserParams.UserDetails.Admin = &falseValue
	UserParams.UserDetails.ProfileUpdatable = &falseValue
	UserParams.UserDetails.DisableUIAccess = &trueValue
	UserParams.UserDetails.InternalPasswordDisabled = &trueValue
	err = testUserService.UpdateUser(UserParams)
	assert.NoError(t, err)
	user, err := testUserService.GetUser(UserParams)

	// We don't know the default group when created, so just set it
	UserParams.UserDetails.Groups = user.Groups
	// Password is not carried in reply
	user.Password = UserParams.UserDetails.Password

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(UserParams.UserDetails, *user))
}

func testClearUserGroups(t *testing.T) {
	UserParams := getTestUserParams(true, "")

	err := testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	UserParams.ClearGroups = true
	err = testUserService.UpdateUser(UserParams)
	assert.NoError(t, err)
	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)

	assert.Nil(t, user.Groups)
}

func testDeleteUser(t *testing.T) {
	UserParams := getTestUserParams(false, "")
	err := testUserService.CreateUser(UserParams)
	assert.NoError(t, err)
	err = testUserService.DeleteUser(UserParams.UserDetails.Name)
	assert.NoError(t, err)
	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func getTestUserParams(replaceIfExists bool, nameSuffix string) services.UserParams {
	userDetails := services.User{
		Name:                     fmt.Sprintf("test%s%s", nameSuffix, timestampStr),
		Email:                    "christianb@jfrog.com",
		Password:                 "Password1*",
		Admin:                    &trueValue,
		Realm:                    "internal",
		ShouldInvite:             &falseValue,
		ProfileUpdatable:         &trueValue,
		DisableUIAccess:          &falseValue,
		InternalPasswordDisabled: &falseValue,
		WatchManager:             &falseValue,
		ReportsManager:           &falseValue,
		PolicyManager:            &falseValue,
	}
	return services.UserParams{
		UserDetails:     userDetails,
		ReplaceIfExists: replaceIfExists,
	}
}

func deleteUserAndAssert(t *testing.T, username string) {
	err := testUserService.DeleteUser(username)
	assert.NoError(t, err)
}
