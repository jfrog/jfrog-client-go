package services

import (
	"net/http"

	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

const batchApi = "api/v2/repositories/batch"

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

func (brs *BatchRepositoryService) PerformBatchRequest(content []byte) (err error) {
	httpClientsDetails := brs.ArtDetails.CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	url := brs.ArtDetails.GetUrl() + batchApi
	var (
		resp *http.Response
		body []byte
	)

	if brs.isUpdate {
		resp, body, err = brs.client.SendPost(url, content, &httpClientsDetails)
	} else {
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

	return nil
}
