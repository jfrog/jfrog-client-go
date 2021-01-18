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
	usersParams := getTestUsersParams(false)

	err := testUsersService.CreateUser(usersParams)
	defer testUsersService.DeleteUser(usersParams.UserDetails.Name)
	assert.NoError(t, err)

	u, _, err := testUsersService.GetUser(usersParams)
	// we don't know the default group when created, so just set it
	usersParams.UserDetails.Groups = u.Groups
	// password is not carried in reply
	u.Password = usersParams.UserDetails.Password
	assert.NotNil(t, u)
	assert.True(t, reflect.DeepEqual(usersParams.UserDetails, *u))

}

func testUpdateUser(t *testing.T) {
	usersParams := getTestUsersParams(true)

	err := testUsersService.CreateUser(usersParams)
	defer testUsersService.DeleteUser(usersParams.UserDetails.Name)
	assert.NoError(t, err)

	usersParams.UserDetails.Email = "changed@mail.com"
	err = testUsersService.UpdateUser(usersParams)
	assert.NoError(t, err)
	user, _, err := testUsersService.GetUser(usersParams)

	// we don't know the default group when created, so just set it
	usersParams.UserDetails.Groups = user.Groups
	// password is not carried in reply
	user.Password = usersParams.UserDetails.Password

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(usersParams.UserDetails, *user))
}

func testDeleteUser(t *testing.T) {
	usersParams := getTestUsersParams(false)
	err := testUsersService.CreateUser(usersParams)
	assert.NoError(t, err)
	err = testUsersService.DeleteUser(usersParams.UserDetails.Name)
	assert.NoError(t, err)
	user, notExists, err := testUsersService.GetUser(usersParams)
	assert.True(t, notExists)
	assert.Nil(t, user)

}

func getTestUsersParams(replaceExistUsers bool) services.UsersParams {
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
	return services.UsersParams{
		UserDetails:       userDetails,
		ReplaceExistUsers: replaceExistUsers,
	}
}
