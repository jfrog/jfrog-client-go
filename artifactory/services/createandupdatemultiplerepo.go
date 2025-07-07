package services

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BatchRepositoryService struct {
	RepositoryService
}

func NewBatchRepositoryService(client *jfroghttpclient.JfrogHttpClient, isUpdate bool) *BatchRepositoryService {
	return &BatchRepositoryService{
		RepositoryService: RepositoryService{
			client:   client,
			isUpdate: isUpdate,
		},
	}
}

func (brs *BatchRepositoryService) PerformBatchRequest(content []byte) error {
	var err error

	httpClientsDetails := brs.ArtDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	url := brs.ArtDetails.GetUrl() + "api/v2/repositories/batch"
	var operationString string
	var resp *http.Response
	var body []byte

	if brs.isUpdate {
		log.Info("Updating multiple repositories...")
		operationString = "updating"
		resp, body, err = brs.client.SendPost(url, content, &httpClientsDetails)
	} else {
		log.Info("Creating multiple repositories...")
		operationString = "creating"
		resp, body, err = brs.client.SendPut(url, content, &httpClientsDetails)
	}

	if err != nil {
		return err
	}
	if brs.isUpdate {
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
			return err
		}
	} else {
		if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusCreated); err != nil {
			return err
		}
	}

	log.Info("Done", operationString, "multiple repositories.")
	return nil
}
