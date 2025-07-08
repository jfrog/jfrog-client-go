package catalog

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-client-go/utils/tests"
	"github.com/stretchr/testify/assert"
)

type MockServerParams struct {
	EnrichedVuln []cyclonedx.Vulnerability
	Alive        bool
}

func StartCatalogMockServerWithParams(t *testing.T, params MockServerParams) int {
	handlers := tests.HttpServerHandlers{}

	handlers["/"] = http.NotFound
	// Version handlers (version is not available in Catalog, so we use ping endpoint)
	handlers["/catalog/api/v1/system/ping"] = catalogGetVersionHandlerFunc(t, params)
	// Enrich handlers
	handlers["/catalog/api/v1/beta/cyclonedx/enrich"] = catalogEnrichHandler(t, params)

	port, err := tests.StartHttpServer(handlers)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return port
}

func catalogGetVersionHandlerFunc(t *testing.T, params MockServerParams) func(w http.ResponseWriter, r *http.Request) {
	version := "1.0.0"
	return func(w http.ResponseWriter, r *http.Request) {
		if !params.Alive {
			http.Error(w, "Catalog service is not available", http.StatusServiceUnavailable)
			return
		}
		if r.Method == http.MethodGet {
			_, err := fmt.Fprint(w, version)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid catalog request", http.StatusBadRequest)
	}
}

func catalogEnrichHandler(t *testing.T, params MockServerParams) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !params.Alive {
			http.Error(w, "Catalog service is not available", http.StatusServiceUnavailable)
			return
		}
		if r.Method == http.MethodPost {
			// Read the BOM from the request body
			bom := cyclonedx.NewBOM()
			assert.NoError(t, cyclonedx.NewBOMDecoder(r.Body, cyclonedx.BOMFileFormatJSON).Decode(bom))
			// Enrich the BOM with vulnerabilities
			for _, vuln := range params.EnrichedVuln {
				if bom.Vulnerabilities == nil {
					bom.Vulnerabilities = &[]cyclonedx.Vulnerability{}
				}
				*bom.Vulnerabilities = append(*bom.Vulnerabilities, vuln)
			}
			// Encode the enriched BOM to JSON format and write it to the response
			writer := cyclonedx.NewBOMEncoder(w, cyclonedx.BOMFileFormatJSON)
			err := writer.Encode(bom)
			assert.NoError(t, err)
			return
		}
		http.Error(w, "Invalid enrich request", http.StatusBadRequest)
	}
}
