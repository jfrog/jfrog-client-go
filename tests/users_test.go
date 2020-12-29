package tests

import (
	"fmt"
	artifactorynew "github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

func TestGroups(t *testing.T) {
	t.Run("create", testCreateUser)
	t.Run("update", testUpdateUser)
	t.Run("delete", testDeleteUser)
}

func testCreateUser(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewUserService(rt.Client())
	gs.SetArtifactoryDetails(details)
	user := services.User{
		Name:                     fmt.Sprintf("test%d", rand.Int()),
		Email:                    "christianb@jfrog.com",
		Password:                 "Password1",
		Admin:                    true,
		Realm:                    "internal",
		ProfileUpdatable:         true,
		DisableUIAccess:          false,
		InternalPasswordDisabled: false,
	}

	err = gs.CreateUser(user)
	assert.NoError(t, err)

	u, err := gs.GetUser(user.Name)
	// we don't know the default group when created, so just set it
	user.Groups = u.Groups
	// password is not carried in reply
	u.Password = user.Password
	assert.NotNil(t, u)
	assert.True(t, reflect.DeepEqual(user, *u))

}

func testUpdateUser(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewUserService(rt.Client())
	gs.SetArtifactoryDetails(details)

	user := services.User{
		Name:                     fmt.Sprintf("test%d", rand.Int()),
		Email:                    "christianb@jfrog.com",
		Password:                 "Password1",
		Admin:                    true,
		Realm:                    "internal",
		ProfileUpdatable:         true,
		DisableUIAccess:          false,
		InternalPasswordDisabled: false,
	}
	err = gs.CreateUser(user)
	assert.NoError(t, err)

	err = gs.UpdateUser(user)
	assert.NoError(t, err)
	usr, err := gs.GetUser(user.Name)

	// we don't know the default group when created, so just set it
	user.Groups = usr.Groups
	// password is not carried in reply
	usr.Password = user.Password

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(user, *usr))
}

func testDeleteUser(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	userService := services.NewUserService(rt.Client())
	userService.SetArtifactoryDetails(details)

	user := services.User{
		Name:                     fmt.Sprintf("test%d", rand.Int()),
		Email:                    "christianb@jfrog.com",
		Password:                 "Password1",
		Admin:                    true,
		Realm:                    "internal",
		ProfileUpdatable:         true,
		DisableUIAccess:          false,
		InternalPasswordDisabled: false,
	}
	err = userService.CreateUser(user)
	assert.NoError(t, err)
	err = userService.DeleteUser(user.Name)
	assert.NoError(t, err)
	g, err := userService.GetUser(user.Name)
	assert.Nil(t, g)

}
