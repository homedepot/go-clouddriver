package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/artifact"
	"github.com/gin-gonic/gin"
)

func ListArtifactCredentials(c *gin.Context) {
	cc := artifact.CredentialsControllerInstance(c)
	c.JSON(http.StatusOK, cc.ListArtifactCredentialsNamesAndTypes())
}

func ListHelmArtifactAccountNames(c *gin.Context) {
	names := []string{}
	cc := artifact.CredentialsControllerInstance(c)

	hc, err := cc.HelmClientForAccountName(c.Param("accountName"))
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	i, err := hc.GetIndex()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	for name := range i.Entries {
		names = append(names, name)
	}

	sort.Strings(names)

	c.JSON(http.StatusOK, names)
}

func ListHelmArtifactAccountVersions(c *gin.Context) {
	cc := artifact.CredentialsControllerInstance(c)
	versions := []string{}
	artifactName := c.Query("artifactName")

	hc, err := cc.HelmClientForAccountName(c.Param("accountName"))
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	i, err := hc.GetIndex()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	found := false

	for name, resources := range i.Entries {
		if strings.EqualFold(name, artifactName) {
			found = true
			for _, resource := range resources {
				versions = append(versions, resource.Version)
			}
		}

		if found {
			break
		}
	}

	sort.Strings(versions)

	c.JSON(http.StatusOK, versions)
}

type Artifact struct {
	Type            artifact.Type `json:"type"`
	CustomKind      bool          `json:"customKind"`
	Name            string        `json:"name"`
	Version         string        `json:"version"`
	Location        string        `json:"location"`
	Reference       string        `json:"reference"`
	Metadata        Metadata      `json:"metadata"`
	ArtifactAccount string        `json:"artifactAccount"`
	// Provenance      interface{} `json:"provenance"`
	// UUID            interface{} `json:"uuid"`
}

type Metadata struct {
	ID string `json:"id"`
}

// This is actually a PUT request to /artifacts/fetch/, but
// I named it "GetArtifact" since that's what it's doing.
func GetArtifact(c *gin.Context) {
	a := Artifact{}
	cc := artifact.CredentialsControllerInstance(c)

	err := c.ShouldBindJSON(&a)
	if err != nil {
		clouddriver.WriteError(c, http.StatusBadRequest, err)
		return
	}

	b := []byte{}

	switch a.Type {
	case artifact.TypeHTTPFile:
		hc, err := cc.HTTPClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		resp, err := hc.Get(a.Reference)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeHelmChart:
		hc, err := cc.HelmClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		b, err = hc.GetChart(a.Name, a.Version)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeEmbeddedBase64:
		// TODO when a base64 encoded helm templated manifest makes its way here, it sometimes starts
		// with "WARNING":"This chart is deprecated" - we should either handle that here by removing
		// the prefix, or handle it in the deploy manifest operation.
		b, err = base64.StdEncoding.DecodeString(a.Reference)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeGithubFile:
		gc, err := cc.GitClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.WriteError(c, http.StatusBadRequest, err)
			return
		}

		if !strings.HasPrefix(a.Reference, gc.BaseURL.String()) {
			clouddriver.WriteError(c, http.StatusBadRequest, fmt.Errorf("content URL %s should have base URL %s",
				a.Reference, gc.BaseURL.String()))
			return
		}

		urlStr := strings.TrimPrefix(a.Reference, gc.BaseURL.String())

		req, err := gc.NewRequest(http.MethodGet, urlStr, nil)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

		branch := "master"
		if a.Version != "" {
			branch = a.Version
		}

		q := req.URL.Query()
		q.Set("ref", branch)
		req.URL.RawQuery = q.Encode()

		var buf bytes.Buffer

		_, err = gc.Do(context.TODO(), req, &buf)
		if err != nil {
			clouddriver.WriteError(c, http.StatusInternalServerError, err)
			return
		}

		var response struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}
		json.Unmarshal(buf.Bytes(), &response)

		if strings.EqualFold(response.Encoding, "base64") {
			b, err = base64.StdEncoding.DecodeString(response.Content)
			if err != nil {
				clouddriver.WriteError(c, http.StatusInternalServerError, err)
				return
			}
		} else {
			b = []byte(response.Content)
		}

	default:
		clouddriver.WriteError(c, http.StatusNotImplemented, fmt.Errorf("getting artifact of type %s not implemented", a.Type))
		return
	}

	c.Writer.Write(b)
}
