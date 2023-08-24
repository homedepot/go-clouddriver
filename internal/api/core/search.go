package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

type SearchResponse []Page

type Page struct {
	PageNumber   int          `json:"pageNumber"`
	PageSize     int          `json:"pageSize"`
	Query        string       `json:"query"`
	Results      []PageResult `json:"results"`
	TotalMatches int          `json:"totalMatches"`
	Platform     string       `json:"platform,omitempty"`
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

var unsupportedKinds = []string{
	"applications",
	"instances",
	"loadBalancers",
	"projects",
	"securityGroups",
}

// Search is the generic search endpoint. It ignores the `pageSize` query parameter
// and lists resources for a kind and namespace across all accounts the user
// has access to concurrently.
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
	// Ignore requests for unsupported kinds
	if containsIgnoreCase(unsupportedKinds, kind) {
		sr := SearchResponse{}
		page := Page{
			PageNumber:   1,
			PageSize:     pageSize,
			Platform:     "aws",
			Query:        namespace,
			Results:      results,
			TotalMatches: 0,
		}
		sr = append(sr, page)

		c.JSON(http.StatusOK, sr)

		return
	}

	wg := &sync.WaitGroup{}
	ac := make(chan accountName, internal.DefaultChanSize)

	providers, err := cc.AllKubernetesProvidersWithTimeout(time.Second * internal.DefaultListTimeoutSeconds)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	providers = filterProviders(providers, accounts)

	wg.Add(len(providers))
	// List the requested resource across accounts concurrently.
	for _, provider := range providers {
		go cc.search(wg, provider, kind, namespace, ac)
	}

	go func() {
		wg.Wait()
		close(ac)
	}()

	// Receive all account/names from the channel.
	accountNames := []accountName{}
	for an := range ac {
		accountNames = append(accountNames, an)
	}

	t := "unclassified"
	if _, ok := spinnakerKindMap[kind]; ok {
		t = spinnakerKindMap[kind]
	}

	for _, an := range accountNames {
		result := PageResult{
			Account:        an.account,
			Group:          kind,
			KubernetesKind: kind,
			Name:           fmt.Sprintf("%s %s", kind, an.name),
			Namespace:      namespace,
			Provider:       "kubernetes",
			Region:         namespace,
			Type:           t,
		}

		results = append(results, result)
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

// accountName defines a Kubernetes resource's account and name.
type accountName struct {
	account string
	name    string
}

// search lists a requested resource by account and kind and namespace
// then writes to a channel of accountName.
func (cc *Controller) search(wg *sync.WaitGroup, provider *kubernetes.Provider,
	kind, namespace string, ac chan accountName) {
	// Increment the wait group counter when we're done here.
	defer wg.Done()

	// If namespace-scoped account and we are attempting to list kinds in a forbidden namespace,
	// just return.
	if err := provider.ValidateNamespaceAccess(namespace); err != nil {
		return
	}
	// If the provider is not allowed to list this given kind, just return.
	if err := provider.ValidateKindStatus(kind); err != nil {
		return
	}
	// Declare a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*internal.DefaultListTimeoutSeconds)
	defer cancel()

	// List resources with the context.
	ul, err := provider.Client.ListResourcesByKindAndNamespaceWithContext(ctx, kind, namespace, metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, u := range ul.Items {
		an := accountName{
			account: provider.Name,
			name:    u.GetName(),
		}
		ac <- an
	}
}

// filterProviders returns a list of providers filtered by the allowe account names passed in.
func filterProviders(providers []*kubernetes.Provider, allowedAccounts []string) []*kubernetes.Provider {
	ps := []*kubernetes.Provider{}
	m := map[string]bool{}

	for _, allowedAccount := range allowedAccounts {
		m[allowedAccount] = true
	}

	for _, provider := range providers {
		if m[provider.Name] {
			ps = append(ps, provider)
		}
	}

	return ps
}
