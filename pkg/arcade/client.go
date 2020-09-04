package arcade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ClientInstanceKey = `ArcadeClient`
)

//go:generate counterfeiter . Client
type Client interface {
	Token() (string, error)
	WithAPIKey(string)
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
	apiKey string
	url    string
}

func (c *client) WithAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *client) Token() (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/tokens", nil)
	if err != nil {
		return "", err
	}

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

func Instance(c *gin.Context) Client {
	return c.MustGet(ClientInstanceKey).(Client)
}
