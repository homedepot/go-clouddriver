package core

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

type SearchResponse []Page

type Page struct {
	PageNumber   int          `json:"pageNumber"`
	PageSize     int          `json:"pageSize"`
	Query        string       `json:"query"`
	Results      []PageResult `json:"results"`
	TotalMatches int          `json:"totalMatches"`
	// Platform     string       `json:"platform"`
}

type PageResult struct {
	Account        string `json:"account"`
	Group          string `json:"group"`
	KubernetesKind string `json:"kubernetesKind"`
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Provider       string `json:"provider"`
	Region         string `json:"region"`
	Type           string `json:"type"`
	Application    string `json:"application,omitempty"`
	Cluster        string `json:"cluster,omitempty"`
}

// Search is the generic search endpoint. It currently just handles searching previously deployed
// Kubernetes resources. It queries the clouddriver DB and does NOT reach out to
// clusters for live data.
func (cc *Controller) Search(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	namespace := c.Query("q")
	// The "type" query param is the kubernetes kind.
	kind := c.Query("type")
	// Get all accounts the user has access to.
	accounts := strings.Split(c.GetHeader("X-Spinnaker-Accounts"), ",")

	if kind == "" || namespace == "" {
		clouddriver.Error(c, http.StatusBadRequest,
			errors.New("must provide query params 'q' to specify the namespace and 'type' to specify the kind"))
		return
	}

	results := []PageResult{}

	for _, account := range accounts {
		if account == "" {
			continue
		}

		if len(results) >= pageSize {
			break
		}

		names, err := cc.SQLClient.ListKubernetesResourceNamesByAccountNameAndKindAndNamespace(account, kind, namespace)
		if err != nil {
			continue
		}

		for _, name := range names {
			t := "unclassified"
			if _, ok := spinnakerKindMap[kind]; ok {
				t = spinnakerKindMap[kind]
			}

			result := PageResult{
				Account:        account,
				Group:          kind,
				KubernetesKind: kind,
				Name:           fmt.Sprintf("%s %s", kind, name),
				Namespace:      namespace,
				Provider:       "kubernetes",
				Region:         namespace,
				Type:           t,
			}

			results = append(results, result)
			if len(results) >= pageSize {
				break
			}
		}
	}

	sr := SearchResponse{}
	page := Page{
		PageNumber:   1,
		PageSize:     pageSize,
		Query:        namespace,
		Results:      results,
		TotalMatches: len(results),
	}
	sr = append(sr, page)

	c.JSON(http.StatusOK, sr)
}
