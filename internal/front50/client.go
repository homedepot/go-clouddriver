package front50

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	defaultFront50Url = "http://spin-front50.spinnaker:8080"
)

//go:generate counterfeiter . Client
type Client interface {
	Project(project string) (Response, error)
}

func NewClient(url string) Client {
	return &client{
		url: url,
	}
}

func NewDefaultClient() Client {
	return NewClient(defaultFront50Url)
}

type client struct {
	url string
}

type Response struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Config         Config `json:"config"`
	UpdateTs       int64  `json:"updateTs"`
	CreateTs       int64  `json:"createTs"`
	LastModifiedBy string `json:"lastModifiedBy"`
}

type Config struct {
	PipelineConfigs []PipelineConfig `json:"pipelineConfigs"`
	Applications    []string         `json:"applications"`
	Clusters        []Cluster        `json:"clusters"`
}

type PipelineConfig struct {
	Application      string `json:"application"`
	PipelineConfigID string `json:"pipelineConfigId"`
}

type Cluster struct {
	Account      string   `json:"account"`
	Stack        string   `json:"stack"`
	Detail       string   `json:"detail"`
	Applications []string `json:"applications"`
}

// Project gets the Spinnaker project from the front50 service
// See https://github.com/spinnaker/front50/blob/master/front50-web/src/main/java/com/netflix/spinnaker/front50/controllers/v2/ProjectsController.java
func (c *client) Project(project string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/v2/projects/"+project, nil)
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
