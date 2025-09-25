package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const (
	catalogEnrichApi              = "api/v1/beta/cyclonedx/enrich"
	catalogMinVersionForEnrichApi = "1.0.0"
)

type EnrichService struct {
	client          *jfroghttpclient.JfrogHttpClient
	CatalogDetails  auth.ServiceDetails
	ScopeProjectKey string
}

func NewEnrichService(client *jfroghttpclient.JfrogHttpClient) *EnrichService {
	return &EnrichService{client: client}
}

func (es *EnrichService) getUrlForEnrichApi() string {
	return utils.AppendScopedProjectKeyParam(es.CatalogDetails.GetUrl()+catalogEnrichApi, es.ScopeProjectKey)
}

// Enrich will enrich the CycloneDX BOM with additional security information
func (es *EnrichService) Enrich(bom *cyclonedx.BOM) (enriched *cyclonedx.BOM, err error) {
	// Encode the BOM to JSON format
	var buf bytes.Buffer
	var writer io.Writer = &buf
	encoder := cyclonedx.NewBOMEncoder(writer, cyclonedx.BOMFileFormatJSON)
	if err = encoder.Encode(bom); err != nil {
		return nil, errorutils.CheckErrorf("failed to encode CycloneDX BOM: %s", err.Error())
	}
	// Enrich the BOM using the Catalog service
	enrichedBom, err := es.enrich(buf.Bytes())
	if err != nil {
		return nil, errorutils.CheckErrorf("failed to enrich CycloneDX BOM: %s", err.Error())
	}
	// Decode the enriched BOM back to a CycloneDX BOM object
	reader := bytes.NewReader(enrichedBom)
	decoder := cyclonedx.NewBOMDecoder(reader, cyclonedx.BOMFileFormatJSON)
	enriched = &cyclonedx.BOM{}
	if err = decoder.Decode(enriched); err != nil {
		return nil, errorutils.CheckErrorf("failed to decode enriched CycloneDX BOM: %s", err.Error())
	}
	return enriched, nil
}

func (es *EnrichService) enrich(bomJson []byte) ([]byte, error) {
	httpDetails := es.CatalogDetails.CreateHttpClientDetails()
	resp, body, err := es.client.SendPost(es.getUrlForEnrichApi(), bomJson, &httpDetails)
	if err != nil {
		return nil, errors.New("failed while attempting to enrich CycloneDX JSON BOM: " + err.Error())
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return nil, errors.New("got unexpected Catalog server response while attempting to enrich CycloneDX JSON BOM:\n" + err.Error())
	}
	return body, nil
}
