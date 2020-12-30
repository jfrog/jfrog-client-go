package tests

import (
	"fmt"
	artifactorynew "github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGroups(t *testing.T) {
	t.Run("create", testCreateGroup)
	t.Run("update", testUpdateGroup)
	t.Run("delete", testDeleteGroup)
}

func testCreateGroup(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewGroupService(rt.Client())
	gs.SetArtifactoryDetails(details)

	group := services.Group{
		Name:            fmt.Sprintf("test%d", rand.Int()),
		Description:     "hello",
		AutoJoin:        false,
		AdminPrivileges: true,
		Realm:           "internal",
		RealmAttributes: "",
	}
	err = gs.CreateOrUpdateGroup(group)
	assert.NoError(t, err)

	g, err := gs.GetGroup(group.Name)
	assert.NotNil(t, g)
	assert.Equal(t, group, *g)
	gs.DeleteGroup(group.Name)

}

func testUpdateGroup(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewGroupService(rt.Client())
	gs.SetArtifactoryDetails(details)

	group := services.Group{
		Name:            fmt.Sprintf("test%d", rand.Int()),
		Description:     "hello",
		AutoJoin:        false,
		AdminPrivileges: true,
		Realm:           "internal",
		RealmAttributes: "",
	}
	err = gs.CreateOrUpdateGroup(group)
	assert.NoError(t, err)

	group.Description = "Changed description"
	err = gs.CreateOrUpdateGroup(group)
	assert.NoError(t, err)
	grp, err := gs.GetGroup(group.Name)
	assert.NoError(t, err)
	assert.Equal(t, group, *grp)
	gs.DeleteGroup(group.Name)

}

func testDeleteGroup(t *testing.T) {
	details := auth.NewArtifactoryDetails()
	details.SetUser("admin")
	details.SetPassword("password")
	details.SetUrl("http://localhost:8081/artifactory/")

	cfg, err := config.NewConfigBuilder().
		SetServiceDetails(details).
		SetDryRun(false).
		Build()
	rt, err := artifactorynew.New(&details, cfg)
	gs := services.NewGroupService(rt.Client())
	gs.SetArtifactoryDetails(details)

	group := services.Group{
		Name:            fmt.Sprintf("test%d", rand.Int()),
		Description:     "hello",
		AutoJoin:        false,
		AdminPrivileges: true,
		Realm:           "internal",
		RealmAttributes: "",
	}
	err = gs.CreateOrUpdateGroup(group)
	assert.NoError(t, err)
	err = gs.DeleteGroup(group.Name)
	assert.NoError(t, err)
	g, err := gs.GetGroup(group.Name)
	assert.Nil(t, g)

}
