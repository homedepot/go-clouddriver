package helm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	ClientInstanceKey = `HelmClient`
)

type Index struct {
	APIVersion string                `json:"apiVersion" yaml:"apiVersion"`
	Entries    map[string][]Resource `json:"entries" yaml:"entries"`
}

type Resource struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	AppVersion string `json:"appVersion" yaml:"appVersion"`
	// Created     time.Time `json:"created"`
	Description string `json:"description" yaml:"description"`
	Digest      string `json:"digest" yaml:"digest"`
	Home        string `json:"home" yaml:"home"`
	// Maintainers []struct {
	// 	Email string `json:"email"`
	// 	Name  string `json:"name"`
	// } `json:"maintainers"`
	Name    string   `json:"name" yaml:"name"`
	Urls    []string `json:"urls" yaml:"urls"`
	Version string   `json:"version" yaml:"version"`
}

//go:generate counterfeiter . Client
type Client interface {
	GetIndex() (Index, error)
	GetChart(string, string) ([]byte, error)
}

var (
	etag  string
	cache Index
	mux   sync.Mutex
)

func NewClient(u string) Client {
	return &client{u: u}
}

type client struct {
	u string
}

func (c *client) GetIndex() (Index, error) {
	i := Index{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/index.yaml", c.u), nil)
	if err != nil {
		return i, err
	}

	req.Header.Add("If-None-Match", etag)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return i, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		mux.Lock()
		defer mux.Unlock()

		return cache, nil
	}

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return i, errors.New("error getting helm index: " + res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return i, err
	}

	err = yaml.Unmarshal(b, &i)
	if err != nil {
		return i, err
	}

	mux.Lock()
	defer mux.Unlock()

	cache = i
	etag = res.Header.Get("etag")

	return i, nil
}

func (c *client) GetChart(name, version string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s-%s.tgz", c.u, name, version))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return nil, errors.New("error getting helm chart: " + res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
