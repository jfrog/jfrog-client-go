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
			_, err := strconv.Atoi(addlSegments[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = fmt.Fprint(w, VulnerabilityReportStatusResponse)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	case http.MethodPost:
		if numSegments == 1 {
			if addlSegments[0] == VulnerabilitiesEndpoint {
				_, err := fmt.Fprint(w, VulnerabilityRequestResponse)
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
			if addlSegments[0] == VulnerabilitiesEndpoint {
				_, err := fmt.Fprint(w, VulnerabilityReportDetailsResponse)
				if err != nil {
					log.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	case http.MethodDelete:
		if numSegments == 0 {
			_, err := fmt.Fprint(w, VulnerabilityReportDeleteResponse)
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

func securityHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	endpoint := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	switch endpoint {
	case "gitinfo":
		_, err = fmt.Fprint(w, gitInfoSentResponse)
	case "graph":
		_, err = fmt.Fprint(w, scanGraphResponse)
	case "9c9dbd61-f544-4e33-4613-34727043d71f":
		_, err = fmt.Fprint(w, getScanResultsResponse)
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

func StartXrayMockServer() int {
	handlers := clienttests.HttpServerHandlers{}
	handlers["/api/xray/scanBuild"] = scanBuildHandler
	handlers["/api/v2/summary/artifact"] = artifactSummaryHandler
	handlers["/api/v1/entitlements/feature/"] = entitlementsHandler
	handlers["/xsc/"] = securityHandler
	handlers["/xray/"] = securityHandler
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
