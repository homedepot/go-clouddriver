package fiat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ClientInstanceKey = `FiatClient`
	fiatUrl           = "http://spin-fiat.spinnaker:7003"
)

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
	Name               string           `json:"name"`
	Accounts           []Account        `json:"accounts"`
	Applications       []Application    `json:"applications"`
	ServiceAccounts    []ServiceAccount `json:"serviceAccounts"`
	Roles              []Role           `json:"roles"`
	BuildServices      []interface{}    `json:"buildServices"`
	ExtensionResources struct {
	} `json:"extensionResources"`
	Admin                            bool `json:"admin"`
	LegacyFallback                   bool `json:"legacyFallback"`
	AllowAccessToUnknownApplications bool `json:"allowAccessToUnknownApplications"`
}

type Account struct {
	Name           string   `json:"name"`
	Authorizations []string `json:"authorizations"`
}

type Application struct {
	Name           string   `json:"name"`
	Authorizations []string `json:"authorizations"`
}

type ServiceAccount struct {
	Name     string   `json:"name"`
	MemberOf []string `json:"memberOf"`
}

type Role struct {
	Name   string `json:"name"`
	Source string `json:"source"`
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

func Instance(c *gin.Context) Client {
	return c.MustGet(ClientInstanceKey).(Client)
}
