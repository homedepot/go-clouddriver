package core

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/artifact"
)

var matchToFirstSlashRegexp = regexp.MustCompile("^[^/]*/")

func ListArtifactCredentials(c *gin.Context) {
	cc := artifact.CredentialsControllerInstance(c)
	c.JSON(http.StatusOK, cc.ListArtifactCredentialsNamesAndTypes())
}

func ListHelmArtifactAccountNames(c *gin.Context) {
	names := []string{}
	cc := artifact.CredentialsControllerInstance(c)

	hc, err := cc.HelmClientForAccountName(c.Param("accountName"))
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	i, err := hc.GetIndex()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	i, err := hc.GetIndex()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	b := []byte{}

	switch a.Type {
	case artifact.TypeHTTPFile:
		hc, err := cc.HTTPClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		resp, err := hc.Get(a.Reference)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeHelmChart:
		hc, err := cc.HelmClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		b, err = hc.GetChart(a.Name, a.Version)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeEmbeddedBase64:
		// TODO when a base64 encoded helm templated manifest makes its way here, it sometimes starts
		// with "WARNING":"This chart is deprecated" - we should either handle that here by removing
		// the prefix, or handle it in the deploy manifest operation.
		b, err = base64.StdEncoding.DecodeString(a.Reference)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeGithubFile:
		gc, err := cc.GitClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		if !strings.HasPrefix(a.Reference, gc.BaseURL.String()) {
			clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("content URL %s should have base URL %s",
				a.Reference, gc.BaseURL.String()))
			return
		}

		urlStr := strings.TrimPrefix(a.Reference, gc.BaseURL.String())

		req, err := gc.NewRequest(http.MethodGet, urlStr, nil)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
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
			clouddriver.Error(c, http.StatusInternalServerError, err)
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
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}
		} else {
			b = []byte(response.Content)
		}

	case artifact.TypeGitRepo:
		gc, err := cc.GitRepoClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		branch := "master"
		if a.Version != "" {
			branch = a.Version
		}

		url := fmt.Sprintf("%s/archive/%s.tar.gz", a.Reference, branch)

		resp, err := gc.Get(url)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			clouddriver.Error(c, http.StatusInternalServerError, fmt.Errorf("error getting git/repo (repo: %s, branch: %s): %s", a.Reference, branch, resp.Status))
			return
		}
		defer resp.Body.Close()

		rb := []byte{}
		rb, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		// The git tarball returns the files in a parent directory that Spinnaker does not expect
		// For example, running this command
		//
		//    curl -sL https://github.com/homedepot/go-clouddriver/archive/master.tar.gz | tar t -
		//
		// Reports this file listing:
		// go-clouddriver-master/
		// go-clouddriver-master/.github/
		// go-clouddriver-master/.github/workflows/
		// go-clouddriver-master/.github/workflows/go.yml
		// go-clouddriver-master/.gitignore
		// go-clouddriver-master/Makefile
		// ...
		//
		// Spinnaker (rosco) expects:
		// .github/
		// .github/workflows/
		// .github/workflows/go.yml
		// .gitignore
		// Makefile
		// ...
		//

		// tar > gzip > buf
		var buf bytes.Buffer
		zr := gzip.NewWriter(&buf)
		tw := tar.NewWriter(zr)

		// Uncompress tarball
		gr, err := gzip.NewReader(bytes.NewBuffer(rb))
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}
		defer gr.Close()

		// Read tar archive
		tr := tar.NewReader(gr)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			tb, err := ioutil.ReadAll(tr)
			if err != nil {
				clouddriver.Error(c, http.StatusInternalServerError, err)
				return
			}

			// Write to new tar archive
			hdr.Name = matchToFirstSlashRegexp.ReplaceAllString(hdr.Name, "")
			if len(hdr.Name) > 0 && strings.HasPrefix(hdr.Name, a.Location) {
				if err := tw.WriteHeader(hdr); err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}

				if _, err := tw.Write(tb); err != nil {
					clouddriver.Error(c, http.StatusInternalServerError, err)
					return
				}
			}
		}

		// Produce tar
		if err := tw.Close(); err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Produce gzip
		if err := zr.Close(); err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		b = buf.Bytes()

	default:
		clouddriver.Error(c, http.StatusNotImplemented, fmt.Errorf("getting artifact of type %s not implemented", a.Type))
		return
	}

	c.Writer.Write(b)
}
