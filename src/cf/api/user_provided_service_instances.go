package api

import (
	"bytes"
	"cf"
	"cf/configuration"
	"cf/net"
	"encoding/json"
	"fmt"
)

type UserProvidedServiceInstanceRepository interface {
	Create(name, drainUrl string, params map[string]string) (apiResponse net.ApiResponse)
	Update(serviceInstanceFields cf.ServiceInstanceFields) (apiResponse net.ApiResponse)
}

type CCUserProvidedServiceInstanceRepository struct {
	config  *configuration.Configuration
	gateway net.Gateway
}

func NewCCUserProvidedServiceInstanceRepository(config *configuration.Configuration, gateway net.Gateway) (repo CCUserProvidedServiceInstanceRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CCUserProvidedServiceInstanceRepository) Create(name, drainUrl string, params map[string]string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/user_provided_service_instances", repo.config.Target)

	type RequestBody struct {
		Name           string            `json:"name"`
		Credentials    map[string]string `json:"credentials"`
		SpaceGuid      string            `json:"space_guid"`
		SysLogDrainUrl string            `json:"syslog_drain_url"`
	}

	jsonBytes, err := json.Marshal(RequestBody{
		Name:           name,
		Credentials:    params,
		SpaceGuid:      repo.config.SpaceFields.Guid,
		SysLogDrainUrl: drainUrl,
	})

	if err != nil {
		apiResponse = net.NewApiResponseWithError("Error parsing response", err)
		return
	}

	return repo.gateway.CreateResource(path, repo.config.AccessToken, bytes.NewReader(jsonBytes))
}

func (repo CCUserProvidedServiceInstanceRepository) Update(serviceInstanceFields cf.ServiceInstanceFields) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/user_provided_service_instances/%s", repo.config.Target, serviceInstanceFields.Guid)

	type RequestBody struct {
		Credentials    map[string]string `json:"credentials,omitempty"`
		SysLogDrainUrl string            `json:"syslog_drain_url,omitempty"`
	}

	reqBody := RequestBody{serviceInstanceFields.Params, serviceInstanceFields.SysLogDrainUrl}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		apiResponse = net.NewApiResponseWithError("Error parsing response", err)
		return
	}

	return repo.gateway.UpdateResource(path, repo.config.AccessToken, bytes.NewReader(jsonBytes))
}
