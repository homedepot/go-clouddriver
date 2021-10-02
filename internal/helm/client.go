package helm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	errUnableToFindResource = errors.New("unable to find resource")
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
	WithUsernameAndPassword(string, string)
	GetChart(string, string) ([]byte, error)
}

var (
	etag  string
	cache Index
	mux   sync.Mutex
)

func NewClient(url string) Client {
	return &client{url: url}
}

type client struct {
	url      string
	username string
	password string
}

func (c *client) WithUsernameAndPassword(username, password string) {
	c.username = username
	c.password = password
}

func (c *client) GetIndex() (Index, error) {
	i := Index{}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/index.yaml", c.url), nil)
	if err != nil {
		return i, err
	}

	req.Header.Add("If-None-Match", etag)

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

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
	var (
		err error
		b   []byte
	)

	resource, err := c.findResource(name, version)
	if err != nil {
		return b, fmt.Errorf("helm: unable to find chart %s-%s: %w", name, version, err)
	}

	if len(resource.Urls) == 0 {
		return b, fmt.Errorf("helm: no resource urls defined for chart %s-%s", name, version)
	}

	// Loop through all the resource's URLs to get the chart.
	for _, url := range resource.Urls {
		req, e := http.NewRequest(http.MethodGet, url, nil)
		if e != nil {
			err = e

			continue
		}
		// Set credentials when chart is hosted in authenticated repository
		if strings.HasPrefix(url, c.url) && c.username != "" && c.password != "" {
			req.SetBasicAuth(c.username, c.password)
		}

		res, e := http.DefaultClient.Do(req)
		if e != nil {
			err = e

			continue
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 399 {
			err = errors.New("helm: error getting chart: " + res.Status)

			continue
		}

		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}

		break
	}

	return b, err
}

// findResource resets the helm index's cache then gets the resource
// from the cache by name and version.
//
// If it is unable to find the resource it returns an error.
func (c *client) findResource(name, version string) (Resource, error) {
	// Refresh the cached index.
	_, err := c.GetIndex()
	if err != nil {
		return Resource{}, err
	}

	// Lock since we are accessing the cached index.
	mux.Lock()
	defer mux.Unlock()

	if _, ok := cache.Entries[name]; ok {
		resources := cache.Entries[name]
		for _, resource := range resources {
			if resource.Version == version {
				return resource, nil
			}
		}
	}

	return Resource{}, errUnableToFindResource
}
