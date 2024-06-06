package core

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
)

var (
	matchToFirstSlashRegexp  = regexp.MustCompile("^[^/]*/")
	matchGcsObjectNameRegexp = regexp.MustCompile(`^gs://(?P<bucket>[^/]*)/(?P<filepath>[^#]*)#?(?P<generation>.*)?`)
)

func (cc *Controller) ListArtifactCredentials(c *gin.Context) {
	c.JSON(http.StatusOK, cc.ArtifactCredentialsController.ListArtifactCredentialsNamesAndTypes())
}

func (cc *Controller) ListHelmArtifactAccountNames(c *gin.Context) {
	names := []string{}

	hc, err := cc.ArtifactCredentialsController.HelmClientForAccountName(c.Param("accountName"))
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

func (cc *Controller) ListHelmArtifactAccountVersions(c *gin.Context) {
	versions := []string{}
	artifactName := c.Query("artifactName")

	hc, err := cc.ArtifactCredentialsController.HelmClientForAccountName(c.Param("accountName"))
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
func (cc *Controller) GetArtifact(c *gin.Context) {
	a := Artifact{}

	err := c.ShouldBindJSON(&a)
	if err != nil {
		clouddriver.Error(c, http.StatusBadRequest, err)
		return
	}

	var b []byte

	switch a.Type {
	case artifact.TypeEmbeddedBase64:
		// TODO when a base64 encoded helm templated manifest makes its way here, it sometimes starts
		// with "WARNING":"This chart is deprecated" - we should either handle that here by removing
		// the prefix, or handle it in the deploy manifest operation.
		b, err = base64.StdEncoding.DecodeString(a.Reference)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeGCSObject:
		// Get GCS client
		gcs, err := cc.ArtifactCredentialsController.GCSClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}
		// Parse filename
		matches := matchGcsObjectNameRegexp.FindStringSubmatch(a.Reference)

		if matches == nil || len(matches) < 2 {
			clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("gcs/object references must be of the format gs://<bucket>/<file-path>[#generation], got: %s", a.Reference))
			return
		}
		// Define GCS object
		object := gcs.Bucket(matches[1]).Object(matches[2])

		if matches[3] != "" {
			generation, err := strconv.ParseInt(matches[3], 10, 64)
			if err != nil {
				clouddriver.Error(c, http.StatusBadRequest, fmt.Errorf("gcs/object generation values must be numeric, got: %s", matches[3]))
				return
			}

			object = object.Generation(generation)
		}
		// Get object reader
		reader, err := object.NewReader(c)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				clouddriver.Error(c, http.StatusNotFound, err)
				return
			}

			clouddriver.Error(c, http.StatusInternalServerError, err)

			return
		}

		defer reader.Close()

		// Read contents
		b, err = io.ReadAll(reader)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeGithubFile:
		gc, err := cc.ArtifactCredentialsController.GitClientForAccountName(a.ArtifactAccount)
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

		_, err = gc.Do(c, req, &buf)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

		var response struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}

		err = json.Unmarshal(buf.Bytes(), &response)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

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
		gc, err := cc.ArtifactCredentialsController.GitRepoClientForAccountName(a.ArtifactAccount)
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
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			clouddriver.Error(c, http.StatusInternalServerError, fmt.Errorf("error getting git/repo (repo: %s, branch: %s): %s", a.Reference, branch, resp.Status))
			return
		}

		rb, err := io.ReadAll(resp.Body)
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

			tb, err := io.ReadAll(tr)
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

	case artifact.TypeHelmChart:
		hc, err := cc.ArtifactCredentialsController.HelmClientForAccountName(a.ArtifactAccount)
		if err != nil {
			clouddriver.Error(c, http.StatusBadRequest, err)
			return
		}

		b, err = hc.GetChart(a.Name, a.Version)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	case artifact.TypeHTTPFile:
		hc, err := cc.ArtifactCredentialsController.HTTPClientForAccountName(a.ArtifactAccount)
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

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			clouddriver.Error(c, http.StatusInternalServerError, err)
			return
		}

	default:
		clouddriver.Error(c, http.StatusNotImplemented, fmt.Errorf("getting artifact of type %s not implemented", a.Type))
		return
	}

	_, err = c.Writer.Write(b)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
	}
}
