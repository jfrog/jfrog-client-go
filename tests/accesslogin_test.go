//go:build itest

package tests

import (
	"encoding/json"
	accessAuth "github.com/jfrog/jfrog-client-go/access/auth"
	"github.com/jfrog/jfrog-client-go/access/services"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	uuid         = "uuid-token"
	accessToken  = "access-token"
	refreshToken = "refresh-token"
)

func TestAccessLogin(t *testing.T) {
	initAccessTest(t)
	getReqNum := 0
	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			assert.Equal(t, "/api/v2/authentication/jfrog_client_login/request", r.URL.Path)

			// Verify body
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			req := services.LoginAuthRequestBody{}
			err = json.Unmarshal(body, &req)
			assert.NoError(t, err)
			assert.Equal(t, uuid, req.Session)

			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			assert.Equal(t, "/api/v2/authentication/jfrog_client_login/token/"+uuid, r.URL.Path)
			getReqNum++
			if getReqNum == 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := auth.CommonTokenParams{AccessToken: accessToken, RefreshToken: refreshToken}
			body, err := json.Marshal(resp)
			assert.NoError(t, err)
			_, err = w.Write(body)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
		default:
			assert.Fail(t, "received unexpected request method")
		}

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createAccessLoginService(t, ts.URL)
	assert.NoError(t, service.SendLoginAuthenticationRequest(uuid))
	token, err := service.GetLoginAuthenticationToken(uuid)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, token.AccessToken)
	assert.Equal(t, refreshToken, token.RefreshToken)
	assert.Equal(t, 2, getReqNum)
}

func TestAccessLoginTimeout(t *testing.T) {
	initAccessTest(t)
	orgMaxWait := services.MaxWait
	defer func() { services.MaxWait = orgMaxWait }()
	services.MaxWait = time.Second

	// Create mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			assert.Equal(t, "/api/v2/authentication/jfrog_client_login/token/"+uuid, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		default:
			assert.Fail(t, "received unexpected request method")
		}

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	service := createAccessLoginService(t, ts.URL)
	_, err := service.GetLoginAuthenticationToken(uuid)
	assert.Error(t, err)
}

func createAccessLoginService(t *testing.T, url string) *services.LoginService {
	rtDetails := accessAuth.NewAccessDetails()
	rtDetails.SetUrl(url + "/")

	// Create http client
	client, err := jfroghttpclient.JfrogClientBuilder().
		SetInsecureTls(true).
		SetClientCertPath(rtDetails.GetClientCertPath()).
		SetClientCertKeyPath(rtDetails.GetClientCertKeyPath()).
		AppendPreRequestInterceptor(rtDetails.RunPreRequestFunctions).
		Build()
	assert.NoError(t, err, "Failed to create JFrog client: %v\n")

	loginService := services.NewLoginService(client)
	loginService.ServiceDetails = rtDetails
	return loginService
}
