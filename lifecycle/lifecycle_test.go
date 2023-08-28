package lifecycle

import (
	"encoding/json"
	artifactoryAuth "github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	lifecycle "github.com/jfrog/jfrog-client-go/lifecycle/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testRb = lifecycle.ReleaseBundleDetails{
	ReleaseBundleName:    "bundle-test",
	ReleaseBundleVersion: "1.2.3",
}

type testCase struct {
	sync        bool
	errExpected bool
	finalStatus lifecycle.RbStatus
}

func TestSimpleGetReleaseBundleStatus(t *testing.T) {
	testCases := map[string]testCase{
		"no sync processing": {sync: false, errExpected: false, finalStatus: lifecycle.Processing},
		"no sync pending":    {sync: false, errExpected: false, finalStatus: lifecycle.Pending},
		"no sync completed":  {sync: false, errExpected: false, finalStatus: lifecycle.Completed},
		"no sync failed":     {sync: false, errExpected: false, finalStatus: lifecycle.Failed},
		"sync completed":     {sync: true, errExpected: false, finalStatus: lifecycle.Completed},
		"sync rejected":      {sync: true, errExpected: false, finalStatus: lifecycle.Rejected},
		"sync failed":        {sync: true, errExpected: false, finalStatus: lifecycle.Failed},
		"sync deleting":      {sync: true, errExpected: false, finalStatus: lifecycle.Deleting},
		"unexpected status":  {sync: true, errExpected: true, finalStatus: "some status"},
	}
	for testName, test := range testCases {
		t.Run(testName, func(t *testing.T) {
			handlerFunc, requestNum := createDefaultHandlerFunc(t, test.finalStatus)
			testGetRBStatus(t, test, handlerFunc)
			assert.Equal(t, 1, *requestNum)
		})
	}

}

func TestComplexReleaseBundleWaitForOperation(t *testing.T) {
	lifecycle.SyncSleepInterval = 1 * time.Second
	defer func() { lifecycle.SyncSleepInterval = lifecycle.DefaultSyncSleepInterval }()

	requestNum := 0
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/"+lifecycle.GetReleaseBundleCreationStatusRestApi(testRb) {
			w.WriteHeader(http.StatusOK)
			var rbStatus lifecycle.RbStatus
			switch requestNum {
			case 0:
				rbStatus = lifecycle.Pending
			case 1:
				rbStatus = lifecycle.Processing
			case 2:
				rbStatus = lifecycle.Completed
			}
			requestNum++
			writeMockStatusResponse(t, w, lifecycle.ReleaseBundleStatusResponse{Status: rbStatus})
		}
	}
	test := testCase{sync: true, errExpected: false, finalStatus: lifecycle.Completed}
	testGetRBStatus(t, test, handlerFunc)
	assert.Equal(t, 3, requestNum)
}

func testGetRBStatus(t *testing.T, test testCase, handlerFunc http.HandlerFunc) {
	mockServer, rbService := createMockServer(t, handlerFunc)
	defer mockServer.Close()

	statusResp, err := rbService.GetReleaseBundleCreationStatus(testRb, "", test.sync)
	if test.errExpected {
		assert.Error(t, err)
		return
	}

	assert.NoError(t, err)
	assert.Equal(t, test.finalStatus, statusResp.Status)
}

func createMockServer(t *testing.T, testHandler http.HandlerFunc) (*httptest.Server, *lifecycle.ReleaseBundlesService) {
	testServer := httptest.NewServer(testHandler)

	rtDetails := artifactoryAuth.NewArtifactoryDetails()
	rtDetails.SetUrl(testServer.URL + "/")

	client, err := jfroghttpclient.JfrogClientBuilder().Build()
	assert.NoError(t, err)
	return testServer, lifecycle.NewReleaseBundlesService(rtDetails, client)
}

func writeMockStatusResponse(t *testing.T, w http.ResponseWriter, statusResp lifecycle.ReleaseBundleStatusResponse) {
	content, err := json.Marshal(statusResp)
	assert.NoError(t, err)
	_, err = w.Write(content)
	assert.NoError(t, err)
}

func createDefaultHandlerFunc(t *testing.T, status lifecycle.RbStatus) (http.HandlerFunc, *int) {
	requestNum := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/"+lifecycle.GetReleaseBundleCreationStatusRestApi(testRb) {
			w.WriteHeader(http.StatusOK)
			requestNum++
			writeMockStatusResponse(t, w, lifecycle.ReleaseBundleStatusResponse{Status: status})
		}
	}, &requestNum
}
