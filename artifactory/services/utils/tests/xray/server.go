package xray

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-client-go/utils/log"
	clienttests "github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/jfrog/jfrog-client-go/xray/services"
	"github.com/stretchr/testify/assert"
)

const (
	CleanScanBuildName          = "cleanBuildName"
	FatalScanBuildName          = "fatalBuildName"
	VulnerableBuildName         = "vulnerableBuildName"
	VulnerabilitiesEndpoint     = "vulnerabilities"
	LicensesEndpoint            = "licenses"
	ContextualAnalysisFeatureId = "contextual_analysis"
	BadFeatureId                = "unknown"
)

func scanBuildHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buildName, err := jsonparser.GetString(body, "buildName")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch buildName {
	case CleanScanBuildName:
		_, err = fmt.Fprint(w, CleanXrayScanResponse)
	case FatalScanBuildName:
		_, err = fmt.Fprint(w, FatalErrorXrayScanResponse)
	case VulnerableBuildName:
		_, err = fmt.Fprint(w, VulnerableXrayScanResponse)
	}
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func artifactSummaryHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprint(w, VulnerableXraySummaryArtifactResponse)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	reportsPathSegmentsCnt := len(strings.Split(services.ReportsAPI, "/"))
	pathSegments := strings.Split(strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/"), "/")
	addlSegments := pathSegments[reportsPathSegmentsCnt:]
	numSegments := len(addlSegments)

	switch r.Method {
	case http.MethodGet:
		if numSegments == 1 {
			id, err := strconv.Atoi(addlSegments[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = fmt.Fprint(w, MapResponse[MapReportIdEndpoint[id]]["ReportStatus"])
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	case http.MethodPost:
		if numSegments == 1 {
			if addlSegments[0] == VulnerabilitiesEndpoint || addlSegments[0] == LicensesEndpoint {
				_, err := fmt.Fprint(w, MapResponse[addlSegments[0]]["XrayReportRequest"])
				if err != nil {
					log.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		} else if numSegments == 2 {
			_, err := strconv.Atoi(addlSegments[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if addlSegments[0] == VulnerabilitiesEndpoint || addlSegments[0] == LicensesEndpoint {
				_, err := fmt.Fprint(w, MapResponse[addlSegments[0]]["ReportDetails"])
				if err != nil {
					log.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	case http.MethodDelete:
		if numSegments == 0 {
			_, err := fmt.Fprint(w, XrayReportDeleteResponse)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "Invalid reports request", http.StatusBadRequest)
}

func entitlementsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	featureId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	switch featureId {
	case ContextualAnalysisFeatureId:
		_, err = fmt.Fprint(w, EntitledResponse)
	case BadFeatureId:
		_, err = fmt.Fprint(w, NotEntitledResponse)
	}
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildScanHandler(w http.ResponseWriter, r *http.Request) {
	argsSegment := strings.Split(r.URL.Path, services.BuildScanAPI)[1]
	switch r.Method {
	case http.MethodGet:
		if argsSegment == "/test-get/3" {
			_, err := fmt.Fprintf(w, BuildScanResultsResponse, "get")
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buildName, err := jsonparser.GetString(body, "build_name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if argsSegment == "/" && strings.HasPrefix(buildName, "test-") {
			_, err = fmt.Fprint(w, TriggerBuildScanResponse)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if argsSegment == "/scanResult" && buildName == "test-post" {
			_, err = fmt.Fprintf(w, BuildScanResultsResponse, "post")
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "Invalid reports request", http.StatusBadRequest)
}

func xscGetVersionHandlerFunc(t *testing.T, version string) func(w http.ResponseWriter, r *http.Request) {
	expectedResponse := fmt.Sprintf(xscVersionResponse, version)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := fmt.Fprint(w, expectedResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid xsc request", http.StatusBadRequest)
	}
}

func xrayGetVersionHandlerFunc(t *testing.T, version string) func(w http.ResponseWriter, r *http.Request) {
	expectedResponse := fmt.Sprintf(xrayVersionResponse, version)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := fmt.Fprint(w, expectedResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid xray request", http.StatusBadRequest)
	}
}

func enrichGetScanId(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_, err := fmt.Fprint(w, scanIdResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid enrich get scan id request", http.StatusBadRequest)
	}
}

func getJasConfig(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := fmt.Fprint(w, JasConfigResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid enrich get scan id request", http.StatusBadRequest)
	}
}

func enrichGetResults(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, err := fmt.Fprint(w, ScanResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid enrich get results request", http.StatusBadRequest)
	}
}

func xscGitInfoHandlerFunc(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		if r.Method == http.MethodPost {
			var reqBody services.XscGitInfoContext
			err = json.Unmarshal(req, &reqBody)
			assert.NoError(t, err)
			if reqBody.GitRepoUrl == "" || reqBody.BranchName == "" || reqBody.CommitHash == "" {
				w.WriteHeader(http.StatusBadRequest)
				_, err := fmt.Fprint(w, XscGitInfoBadResponse)
				assert.NoError(t, err)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_, err = fmt.Fprint(w, XscGitInfoResponse)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid xsc request", http.StatusBadRequest)
	}
}

type MockServerParams struct {
	MSI         string
	XrayVersion string
	XscVersion  string
}

func StartXrayMockServer(t *testing.T) int {
	params := MockServerParams{MSI: TestMultiScanId, XrayVersion: "3.0.0", XscVersion: "1.0.0"}
	return StartXrayMockServerWithParams(t, params)
}

func StartXrayMockServerWithParams(t *testing.T, params MockServerParams) int {
	handlers := clienttests.HttpServerHandlers{}

	handlers["/xsc/api/v1/gitinfo"] = xscGitInfoHandlerFunc(t)

	handlers["/"] = http.NotFound
	// Xray handlers
	handlers["/xray/api/v1/system/version"] = xrayGetVersionHandlerFunc(t, params.XrayVersion)
	handlers["/api/xray/scanBuild"] = scanBuildHandler
	handlers["/api/v2/summary/artifact"] = artifactSummaryHandler
	handlers["/api/v1/entitlements/feature/"] = entitlementsHandler
	handlers["/xray/api/v1/scan/import_xml"] = enrichGetScanId(t)
	handlers[fmt.Sprintf("/xray/api/v1/scan/graph/%s", params.MSI)] = enrichGetResults(t)
	handlers["/xray/api/v1/configuration/jas"] = getJasConfig(t)
	handlers[fmt.Sprintf("/%s/", services.BuildScanAPI)] = buildScanHandler
	handlers[fmt.Sprintf("/%s/", services.ReportsAPI)] = reportHandler
	// Xsc handlers
	handlers["/xsc/api/v1/system/version"] = xscGetVersionHandlerFunc(t, params.XscVersion)
	handlers["/xray/api/v1/xsc/api/v1/system/version"] = xscGetVersionHandlerFunc(t, params.XscVersion)

	port, err := clienttests.StartHttpServer(handlers)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return port
}
