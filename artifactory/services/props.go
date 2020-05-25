package services

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/jfrog/gofrog/parallel"
	rthttpclient "github.com/jfrog/jfrog-client-go/artifactory/httpclient"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/auth"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type PropsService struct {
	client     *rthttpclient.ArtifactoryHttpClient
	ArtDetails auth.ServiceDetails
	Threads    int
}

func NewPropsService(client *rthttpclient.ArtifactoryHttpClient) *PropsService {
	return &PropsService{client: client}
}

func (ps *PropsService) GetArtifactoryDetails() auth.ServiceDetails {
	return ps.ArtDetails
}

func (ps *PropsService) SetArtifactoryDetails(rt auth.ServiceDetails) {
	ps.ArtDetails = rt
}

func (ps *PropsService) IsDryRun() bool {
	return false
}

func (ps *PropsService) GetThreads() int {
	return ps.Threads
}

func (ps *PropsService) SetProps(propsParams PropsParams) (int, error) {
	log.Info("Setting properties...")
	totalSuccess, err := ps.performRequest(propsParams, false)
	if err != nil {
		return totalSuccess, err
	}
	log.Info("Done setting properties.")
	return totalSuccess, nil
}

func (ps *PropsService) DeleteProps(propsParams PropsParams) (int, error) {
	log.Info("Deleting properties...")
	totalSuccess, err := ps.performRequest(propsParams, true)
	if err != nil {
		return totalSuccess, err
	}
	log.Info("Done deleting properties.")
	return totalSuccess, nil
}

type PropsParams struct {
	Items []utils.ResultItem
	Props string
}

func (sp *PropsParams) GetItems() []utils.ResultItem {
	return sp.Items
}

func (sp *PropsParams) GetProps() string {
	return sp.Props
}

func (ps *PropsService) performRequest(propsParams PropsParams, isDelete bool) (int, error) {
	updatePropertiesBaseUrl := ps.GetArtifactoryDetails().GetUrl() + "api/storage"
	var encodedParam string
	if !isDelete {
		props, err := utils.ParseProperties(propsParams.GetProps(), utils.JoinCommas)
		if err != nil {
			return 0, err
		}
		encodedParam = props.ToEncodedString()
	} else {
		propList := strings.Split(propsParams.GetProps(), ",")
		for _, prop := range propList {
			encodedParam += url.QueryEscape(prop) + ","
		}
		// Remove trailing comma
		if strings.HasSuffix(encodedParam, ",") {
			encodedParam = encodedParam[:len(encodedParam)-1]
		}

	}

	successCounters := make([]int, ps.GetThreads())
	producerConsumer := parallel.NewBounedRunner(ps.GetThreads(), true)
	errorsQueue := clientutils.NewErrorsQueue(1)

	go func() {
		for _, item := range propsParams.GetItems() {
			relativePath := item.GetItemRelativePath()
			setPropsTask := func(threadId int) error {
				logMsgPrefix := clientutils.GetLogMsgPrefix(threadId, ps.IsDryRun())
				setPropertiesUrl := updatePropertiesBaseUrl + "/" + relativePath + "?properties=" + encodedParam
				var resp *http.Response
				var err error
				if isDelete {
					resp, _, err = ps.sendDeleteRequest(logMsgPrefix, relativePath, setPropertiesUrl)
				} else {
					resp, _, err = ps.sendPutRequest(logMsgPrefix, relativePath, setPropertiesUrl)
				}

				if err != nil {
					return err
				}
				if err = errorutils.CheckResponseStatus(resp, http.StatusNoContent); err != nil {
					return errorutils.CheckError(err)
				}
				successCounters[threadId]++
				return nil
			}

			producerConsumer.AddTaskWithError(setPropsTask, errorsQueue.AddError)
		}
		defer producerConsumer.Done()
	}()

	producerConsumer.Run()
	totalSuccess := 0
	for _, v := range successCounters {
		totalSuccess += v
	}

	err := errorsQueue.GetError()
	if err != nil {
		return totalSuccess, err
	}
	return totalSuccess, nil
}

func (ps *PropsService) sendDeleteRequest(logMsgPrefix, relativePath, setPropertiesUrl string) (resp *http.Response, body []byte, err error) {
	log.Info(logMsgPrefix+"Deleting properties on:", relativePath)
	log.Debug(logMsgPrefix+"Sending delete properties request:", setPropertiesUrl)
	httpClientsDetails := ps.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err = ps.client.SendDelete(setPropertiesUrl, nil, &httpClientsDetails)
	return
}

func (ps *PropsService) sendPutRequest(logMsgPrefix, relativePath, setPropertiesUrl string) (resp *http.Response, body []byte, err error) {
	log.Info(logMsgPrefix+"Setting properties on:", relativePath)
	log.Debug(logMsgPrefix+"Sending set properties request:", setPropertiesUrl)
	httpClientsDetails := ps.GetArtifactoryDetails().CreateHttpClientDetails()
	resp, body, err = ps.client.SendPut(setPropertiesUrl, nil, &httpClientsDetails)
	return
}

func NewPropsParams() PropsParams {
	return PropsParams{}
}
