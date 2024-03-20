package services

import (
	"encoding/json"
	"net/http"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type ExportService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails auth.ServiceDetails
	// If true, the export will only print the parameters
	DryRun bool
}

func NewExportService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *ExportService {
	return &ExportService{artDetails: artDetails, client: client}
}

func (drs *ExportService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return drs.client
}

func (drs *ExportService) Export(exportParams ExportParams) error {
	httpClientsDetails := drs.artDetails.CreateHttpClientDetails()
	requestContent, err := json.Marshal(ExportBody(exportParams))
	if err != nil {
		return errorutils.CheckError(err)
	}

	exportMessage := "Running full system export..."
	if drs.DryRun {
		log.Info("[Dry run] " + exportMessage)
		log.Info("Export parameters: \n" + clientutils.IndentJson(requestContent))
		return nil
	}
	log.Info(exportMessage)

	utils.SetContentType("application/json", &httpClientsDetails.Headers)
	resp, body, err := drs.client.SendPost(drs.artDetails.GetUrl()+"api/export/system", requestContent, &httpClientsDetails)
	if err != nil {
		return err
	}
	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	log.Info(string(body))
	log.Debug("Artifactory response:", resp.Status)
	return nil
}

type ExportParams struct {
	// Mandatory:
	// A path to a directory on the local file system of Artifactory server
	ExportPath string

	// Optional:
	// If true, repository metadata is included in export (Maven 2 metadata is unaffected by this setting)
	IncludeMetadata *bool
	// If true, creates and exports to a Zip archive
	CreateArchive *bool
	// If true, prints more verbose logging
	Verbose *bool
	// If true, includes Maven 2 repository metadata and checksum files as part of the export
	M2 *bool
	// If true, repository binaries are excluded from the export
	ExcludeContent *bool
}

type ExportBody struct {
	ExportPath      string `json:"exportPath,omitempty"`
	IncludeMetadata *bool  `json:"includeMetadata,omitempty"`
	CreateArchive   *bool  `json:"createArchive,omitempty"`
	Verbose         *bool  `json:"verbose,omitempty"`
	M2              *bool  `json:"m2,omitempty"`
	ExcludeContent  *bool  `json:"excludeContent,omitempty"`
}

func NewExportParams(exportPath string) ExportParams {
	return ExportParams{ExportPath: exportPath}
}
