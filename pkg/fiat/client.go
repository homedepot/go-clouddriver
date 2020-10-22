package fiat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const fiatUrl = "http://spin-fiat.spinnaker:7003"

//go:generate counterfeiter . Client
type Client interface {
	Authorize(account string) (Response, error)
}

func NewClient(url string) Client {
	return &client{
		url: url,
	}
}

func NewDefaultClient() Client {
	return NewClient(fiatUrl)
}

type client struct {
	url string
}

type Response struct {
	Name     string `json:"name"`
	Accounts []struct {
		Name           string   `json:"name"`
		Authorizations []string `json:"authorizations"`
	} `json:"accounts"`
	Applications []struct {
		Name           string   `json:"name"`
		Authorizations []string `json:"authorizations"`
	} `json:"applications"`
	ServiceAccounts []struct {
		Name     string   `json:"name"`
		MemberOf []string `json:"memberOf"`
	} `json:"serviceAccounts"`
	Roles []struct {
		Name   string `json:"name"`
		Source string `json:"source"`
	} `json:"roles"`
	BuildServices      []interface{} `json:"buildServices"`
	ExtensionResources struct {
	} `json:"extensionResources"`
	Admin                            bool `json:"admin"`
	LegacyFallback                   bool `json:"legacyFallback"`
	AllowAccessToUnknownApplications bool `json:"allowAccessToUnknownApplications"`
}

func (c *client) Authorize(account string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/authorize/"+account, nil)
	if err != nil {
		return Response{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return Response{}, fmt.Errorf("user authorization error: %s", res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	response := Response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
