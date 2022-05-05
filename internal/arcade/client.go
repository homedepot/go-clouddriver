package arcade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//go:generate counterfeiter . Client
type Client interface {
	Token(string) (string, error)
	WithAPIKey(string)
	WithShortExpiration(int)
}

func NewDefaultClient() Client {
	return &client{
		url: "http://localhost:1982",
	}
}

func NewClient(url string) Client {
	return &client{
		url: url,
	}
}

type client struct {
	apiKey          string
	shortExpiration int
	url             string
}

func (c *client) WithAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *client) WithShortExpiration(seconds int) {
	c.shortExpiration = seconds
}

func (c *client) Token(tokenProvider string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/tokens", nil)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("provider", tokenProvider)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Api-Key", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return "", fmt.Errorf("error getting token: %s", res.Status)
	}

	var response struct {
		Token string `json:"token"`
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}
