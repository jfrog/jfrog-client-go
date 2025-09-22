package tests

import (
	"fmt"
	"reflect"
	"testing"

	services "github.com/jfrog/jfrog-client-go/access/services/v2"
	"github.com/stretchr/testify/assert"
)

func TestAccessUsers(t *testing.T) {
	initAccessTest(t)
	t.Run("create-update-delete", testAccessUserCreateUpdateDelete)
}

func getTestAccessUserParams(nameSuffix string) services.UserParams {
	userDetails := services.CommonUserParams{
		Username:                 fmt.Sprintf("test%s%s", nameSuffix, timestampStr),
		Email:                    "john.doe@example.com",
		Password:                 "Password1*",
		Admin:                    &trueValue,
		Realm:                    "internal",
		ProfileUpdatable:         &trueValue,
		DisableUIAccess:          &falseValue,
		InternalPasswordDisabled: &falseValue,
	}
	return services.UserParams{
		CommonUserParams: userDetails,
	}
}

func testAccessUserCreateUpdateDelete(t *testing.T) {
	userParams := getTestAccessUserParams("testProject")
	assert.NoError(t, testAccessUserService.Create(userParams))
	defer deleteAccessUserAndAssert(t, userParams.Username)
	userParams.Email = "joe.blow@example.org"
	userParams.Admin = &trueValue
	userParams.ProfileUpdatable = &falseValue
	userParams.DisableUIAccess = &trueValue
	userParams.InternalPasswordDisabled = &trueValue

	assert.NoError(t, testAccessUserService.Update(userParams))

	updatedUser, err := testAccessUserService.Get(userParams.Username)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser, "Expected user %s but got nil", userParams.Username)
	assert.Equal(t, "active", updatedUser.Status, "Expected user in status 'active' but got status '%s'", updatedUser.Status)
	if !reflect.DeepEqual(userParams, updatedUser.CommonUserParams) {
		t.Error("Unexpected user details built. Expected: `", userParams, "` Got `", updatedUser.CommonUserParams, "`")
	}
}

func deleteAccessUserAndAssert(t *testing.T, username string) {
	assert.NoError(t, testAccessUserService.Delete(username))
}

func TestGetAllAccessUsers(t *testing.T) {
	initAccessTest(t)
	preUsers, err := testAccessUserService.GetAll()
	assert.NoError(t, err)

	noOfUsersBeforeTest := len(preUsers)
	user1 := getTestAccessUserParams("u1")
	user2 := getTestAccessUserParams("u2")
	user3 := getTestAccessUserParams("u3")
	user4 := getTestAccessUserParams("u4")

	assert.NoError(t, testAccessUserService.Create(user1))
	assert.NoError(t, testAccessUserService.Create(user2))
	assert.NoError(t, testAccessUserService.Create(user3))
	assert.NoError(t, testAccessUserService.Create(user4))

	projects, err := testAccessUserService.GetAll()
	assert.NoError(t, err, "Failed to Unmarshal")
	assert.Equal(t, noOfUsersBeforeTest+4, len(projects))

	deleteUserAndAssert(t, user1.Username)
	deleteUserAndAssert(t, user2.Username)
	deleteUserAndAssert(t, user3.Username)
	deleteUserAndAssert(t, user4.Username)
}
