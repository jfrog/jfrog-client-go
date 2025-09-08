//go:build itest

package tests

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"net/http"
	"net/http/httptest"
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
	t.Run("get locked users", testGetLockedUsers)
	t.Run("unlock user", testUnlockUser)
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
	UserParams.UserDetails.Admin = utils.Pointer(false)
	UserParams.UserDetails.ProfileUpdatable = utils.Pointer(false)
	UserParams.UserDetails.DisableUIAccess = utils.Pointer(true)
	UserParams.UserDetails.InternalPasswordDisabled = utils.Pointer(true)
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
		Admin:                    utils.Pointer(true),
		Realm:                    "internal",
		ShouldInvite:             utils.Pointer(false),
		ProfileUpdatable:         utils.Pointer(true),
		DisableUIAccess:          utils.Pointer(false),
		InternalPasswordDisabled: utils.Pointer(false),
		WatchManager:             utils.Pointer(false),
		ReportsManager:           utils.Pointer(false),
		PolicyManager:            utils.Pointer(false),
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

func testGetLockedUsers(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		assert.Equal(t, http.MethodGet, r.Method)

		// Check URL
		assert.Equal(t, "/api/security/lockedUsers", r.URL.Path)

		// Send response 200 OK
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("[ \"froguser\" ]"))
		assert.NoError(t, err)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createMockUserService(t, ts.URL)
	results, err := service.GetLockedUsers()
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "froguser", results[0])
}

func testUnlockUser(t *testing.T) {
	err := testUserService.UnlockUser("froguser")
	assert.NoError(t, err)
}

func createMockUserService(t *testing.T, url string) *services.UserService {
	// Create artifactory details
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(url + "/")

	// Create http client
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err, "Failed to create Artifactory client: %v\n")

	// Create system service
	userService := services.NewUserService(client)
	userService.SetArtifactoryDetails(rtDetails)
	return userService
}
