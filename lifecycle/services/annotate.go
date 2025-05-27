package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/distribution"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"path"
	"strings"
)

const (
	Tag               = "tag"
	PropertiesBaseApi = "api/storage"
)

type CommonPropParams struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

type ArtCommonParams struct {
	Url string `json:"url"`
}

type AnnotateOperationParams struct {
	RbTag          RbAnnotationTag
	RbProps        RbAnnotationProps
	RbDelProps     RbDelProps
	RbDetails      ReleaseBundleDetails
	QueryParams    CommonOptionalQueryParams
	PropertyParams CommonPropParams
	ArtifactoryUrl ArtCommonParams
}

type RbDelProps struct {
	Keys  string `json:"keys"`
	Exist bool   `json:"exist"`
}

type RbAnnotationTag struct {
	Tag   string `json:"tag,omitempty"`
	Exist bool   `json:"exist"`
}

type RbAnnotationProps struct {
	Properties map[string][]string `json:"properties,omitempty"`
	Exist      bool                `json:"exist"`
}

func (rbs *ReleaseBundlesService) AnnotateReleaseBundle(params AnnotateOperationParams) error {
	return rbs.annotateReleaseBundle(params)
}

func GetReleaseBundleSetTagApi(rbDetails ReleaseBundleDetails) string {
	return path.Join(releaseBundleBaseApi, records, rbDetails.ReleaseBundleName, rbDetails.ReleaseBundleVersion, Tag)
}

func (rbs *ReleaseBundlesService) setTag(params AnnotateOperationParams, api string, httpClientsDetails httputils.HttpClientDetails) error {
	if !params.RbTag.Exist {
		log.Debug("Tag doesn't exist or empty")
		return nil
	}

	projParam := distribution.GetProjectQueryParam(params.QueryParams.ProjectKey)
	tagContent, err := json.Marshal(params.RbTag)
	if err != nil {
		return errorutils.CheckError(err)
	}

	setTagFullUrl, err := utils.BuildUrl(rbs.GetLifecycleDetails().GetUrl(), api, projParam)
	if err != nil {
		return err
	}

	resp, body, err := rbs.client.SendPut(setTagFullUrl, tagContent, &httpClientsDetails)
	if err != nil {
		return err
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, body, http.StatusOK); err != nil {
		return err
	}
	return err
}

func (rbs *ReleaseBundlesService) setProperties(params AnnotateOperationParams, api string, httpClientsDetails httputils.HttpClientDetails) error {
	if !params.RbProps.Exist {
		log.Debug("Properties doesn't exist or empty")
		return nil
	}

	err := rbs.sendPropertiesRequest(http.MethodPut, params.ArtifactoryUrl.Url,
		api+"/"+params.PropertyParams.Path, GetPropQueryParam(params), httpClientsDetails)
	if err != nil {
		return errorutils.CheckError(err)
	}
	return nil
}

func (rbs *ReleaseBundlesService) deleteProperties(params AnnotateOperationParams, api string, httpClientsDetails httputils.HttpClientDetails) error {
	if !params.RbDelProps.Exist {
		log.Debug("Delete Properties doesn't exist or empty")
		return nil
	}

	err := rbs.sendPropertiesRequest(http.MethodDelete, params.ArtifactoryUrl.Url,
		api+"/"+params.PropertyParams.Path, GetDeletePropParams(params), httpClientsDetails)
	if err != nil {
		return errorutils.CheckError(err)
	}
	return nil
}

func (rbs *ReleaseBundlesService) sendPropertiesRequest(method, url, path string, params map[string]string,
	httpClientsDetails httputils.HttpClientDetails) error {
	var propsFullUrl string
	var err error
	propsFullUrl, err = utils.BuildUrl(url, path, params)
	if err != nil {
		return errorutils.CheckError(err)
	}

	var resp *http.Response
	var respBody []byte
	switch method {
	case http.MethodDelete:
		resp, respBody, err = rbs.client.SendDelete(propsFullUrl, nil, &httpClientsDetails)
		if err != nil {
			return errorutils.CheckError(err)
		}
	case http.MethodPut:
		resp, respBody, err = rbs.client.SendPut(propsFullUrl, nil, &httpClientsDetails)
		if err != nil {
			return errorutils.CheckError(err)
		}
	default:
		return errors.New("Unexpected method: " + method)
	}

	if err = errorutils.CheckResponseStatusWithBody(resp, respBody, http.StatusNoContent, http.StatusOK); err != nil {
		return errorutils.CheckError(err)
	}
	return nil
}

func GetDeletePropParams(params AnnotateOperationParams) map[string]string {
	queryParams := make(map[string]string)
	queryParams["properties"] = params.RbDelProps.Keys
	if params.PropertyParams.Recursive {
		queryParams["recursive"] = "1"
	} else {
		queryParams["recursive"] = "0"
	}

	return queryParams
}

func (rbs *ReleaseBundlesService) annotateReleaseBundle(params AnnotateOperationParams) error {

	httpClientsDetails := rbs.GetLifecycleDetails().CreateHttpClientDetails()
	httpClientsDetails.SetContentTypeApplicationJson()

	err := rbs.setTag(params, GetReleaseBundleSetTagApi(params.RbDetails), httpClientsDetails)
	if err != nil {
		log.Info("Failed to set tag: " + params.RbTag.Tag)
		return err
	}

	err = rbs.setProperties(params, PropertiesBaseApi, httpClientsDetails)
	if err != nil {
		log.Info("Failed to set properties")
		return err
	}

	err = rbs.deleteProperties(params, PropertiesBaseApi, httpClientsDetails)
	if err != nil {
		log.Info("Failed to delete properties")
		return err
	}

	return nil
}

func GetPropQueryParam(params AnnotateOperationParams) map[string]string {
	queryParams := make(map[string]string)
	queryParams["properties"] = params.RbProps.PropertiesToString()
	if params.PropertyParams.Recursive {
		queryParams["recursive"] = "1"
	} else {
		queryParams["recursive"] = "0"
	}

	return queryParams
}
func (r *RbAnnotationProps) PropertiesToString() string {
	var parts []string
	for key, values := range r.Properties {
		parts = append(parts, fmt.Sprintf("%s=%s", key, strings.Join(values, ",")))
	}
	return strings.Join(parts, ";")
}
