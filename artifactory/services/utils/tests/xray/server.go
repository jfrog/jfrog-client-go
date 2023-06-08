package xray

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/jfrog/jfrog-client-go/utils/log"
	clienttests "github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/jfrog/jfrog-client-go/xray/services"
)

const (
	CleanScanBuildName          = "cleanBuildName"
	FatalScanBuildName          = "fatalBuildName"
	VulnerableBuildName         = "vulnerableBuildName"
	VulnerabilitiesEndpoint     = "vulnerabilities"
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
		fmt.Fprint(w, CleanXrayScanResponse)
		return
	case FatalScanBuildName:
		fmt.Fprint(w, FatalErrorXrayScanResponse)
		return
	case VulnerableBuildName:
		fmt.Fprint(w, VulnerableXrayScanResponse)
		return
	}
}

func artifactSummaryHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, VulnerableXraySummaryArtifactResponse)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	reportsPathSegmentsCnt := len(strings.Split(services.ReportsAPI, "/"))
	pathSegments := strings.Split(strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/"), "/")
	addlSegments := pathSegments[reportsPathSegmentsCnt:]
	numSegments := len(addlSegments)

	switch r.Method {
	case http.MethodGet:
		if numSegments == 1 {
			_, err := strconv.Atoi(addlSegments[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, VulnerabilityReportStatusResponse)
			return
		}
	case http.MethodPost:
		if numSegments == 1 {
			if addlSegments[0] == VulnerabilitiesEndpoint {
				_, _ = fmt.Fprint(w, VulnerabilityRequestResponse)
				return
			}
		} else if numSegments == 2 {
			_, err := strconv.Atoi(addlSegments[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if addlSegments[0] == VulnerabilitiesEndpoint {
				_, _ = fmt.Fprint(w, VulnerabilityReportDetailsResponse)
				return
			}
		}
	case http.MethodDelete:
		if numSegments == 0 {
			_, _ = fmt.Fprint(w, VulnerabilityReportDeleteResponse)
			return
		}
	}
	http.Error(w, "Invalid reports request", http.StatusBadRequest)
}

func entitlementsHandler(w http.ResponseWriter, r *http.Request) {
	featureId := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	switch featureId {
	case ContextualAnalysisFeatureId:
		_, _ = fmt.Fprint(w, EntitledResponse)
		return
	case BadFeatureId:
		_, _ = fmt.Fprint(w, NotEntitledResponse)
		return
	}
}

func buildScanHandler(w http.ResponseWriter, r *http.Request) {
	argsSegment := strings.Split(r.URL.Path, services.BuildScanAPI)[1]
	switch r.Method {
	case http.MethodGet:
		if argsSegment == "/test-get/3" {
			_, _ = fmt.Fprintf(w, BuildScanResultsResponse, "get")
			return
		}
	case http.MethodPost:
		// TODO: validate body
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
			_, _ = fmt.Fprint(w, TriggerBuildScanResponse)
			return
		}
		if argsSegment == "/scanResult" && buildName == "test-post" {
			_, _ = fmt.Fprintf(w, BuildScanResultsResponse, "post")
			return
		}
	}
	http.Error(w, "Invalid reports request", http.StatusBadRequest)
}

func StartXrayMockServer() int {
	handlers := clienttests.HttpServerHandlers{}
	handlers["/api/xray/scanBuild"] = scanBuildHandler
	handlers["/api/v2/summary/artifact"] = artifactSummaryHandler
	handlers["/api/v1/entitlements/feature/"] = entitlementsHandler
	handlers[fmt.Sprintf("/%s/", services.ReportsAPI)] = reportHandler
	handlers[fmt.Sprintf("/%s/", services.BuildScanAPI)] = buildScanHandler
	handlers["/"] = http.NotFound

	port, err := clienttests.StartHttpServer(handlers)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return port
}
