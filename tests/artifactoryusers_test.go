package tests

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	t.Run("create", testCreateUser)
	t.Run("update", testUpdateUser)
	t.Run("delete", testDeleteUser)
}

func testCreateUser(t *testing.T) {
	UserParams := getTestUserParams(false)

	err := testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	user, err := testUserService.GetUser(UserParams)
	// we don't know the default group when created, so just set it
	UserParams.UserDetails.Groups = user.Groups
	// password is not carried in reply
	user.Password = UserParams.UserDetails.Password
	assert.NotNil(t, user)
	assert.True(t, reflect.DeepEqual(UserParams.UserDetails, *user))
}

func testUpdateUser(t *testing.T) {
	UserParams := getTestUserParams(true)

	err := testUserService.CreateUser(UserParams)
	defer deleteUserAndAssert(t, UserParams.UserDetails.Name)
	assert.NoError(t, err)

	UserParams.UserDetails.Email = "changed@mail.com"
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

func testDeleteUser(t *testing.T) {
	UserParams := getTestUserParams(false)
	err := testUserService.CreateUser(UserParams)
	assert.NoError(t, err)
	err = testUserService.DeleteUser(UserParams.UserDetails.Name)
	assert.NoError(t, err)
	user, err := testUserService.GetUser(UserParams)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func getTestUserParams(replaceIfExists bool) services.UserParams {
	userDetails := services.User{
		Name:                     fmt.Sprintf("test%d", rand.Int()),
		Email:                    "christianb@jfrog.com",
		Password:                 "Password1",
		Admin:                    false,
		Realm:                    "internal",
		ProfileUpdatable:         true,
		DisableUIAccess:          false,
		InternalPasswordDisabled: false,
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
