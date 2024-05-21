package artifact

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/v32/github"
	"github.com/homedepot/go-clouddriver/internal/helm"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type Type string

const (
	TypeHelmChart               Type = "helm/chart"
	TypeGitRepo                 Type = "git/repo"
	TypeFront50PipelineTemplate Type = "front50/pipelineTemplate"
	TypeEmbeddedBase64          Type = "embedded/base64"
	TypeCustomerObject          Type = "custom/object"
	TypeGCSObject               Type = "gcs/object"
	TypeHTTPFile                Type = "http/file"
	TypeDockerImage             Type = "docker/image"
	TypeKubernetesConfigMap     Type = "kubernetes/configMap"
	TypeKubernetesDeployment    Type = "kubernetes/deployment"
	TypeKubernetesReplicaSet    Type = "kubernetes/replicaSet"
	TypeKubernetesSecret        Type = "kubernetes/secret"
	TypeGithubFile              Type = "github/file"
)

//go:generate counterfeiter . CredentialsController
type CredentialsController interface {
	ListArtifactCredentialsNamesAndTypes() []Credentials
	HelmClientForAccountName(string) (helm.Client, error)
	HTTPClientForAccountName(string) (*http.Client, error)
	GCSClientForAccountName(string) (*storage.Client, error)
	GitClientForAccountName(string) (*github.Client, error)
	GitRepoClientForAccountName(string) (*http.Client, error)
}

type Credentials struct {
	// General config.
	Name  string `json:"name"`
	Types []Type `json:"types"`
	// Helm repository config.
	Repository string `json:"repository,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	// Github config.
	BaseURL    string `json:"baseURL,omitempty"`
	Token      string `json:"token,omitempty"`
	Enterprise bool   `json:"enterprise,omitempty"`
	// GCS Object config.
	JSONPath string `json:"jsonPath,omitempty"`
}

const (
	defaultConfigDir = "/opt/spinnaker/artifacts/config"
)

func NewDefaultCredentialsController() (CredentialsController, error) {
	return NewCredentialsController(defaultConfigDir)
}

func NewCredentialsController(dir string) (CredentialsController, error) {
	cc := credentialsController{
		artifactCredentials: []Credentials{},
		gcsClients:          map[string]*storage.Client{},
		gitClients:          map[string]*github.Client{},
		gitRepoClients:      map[string]*http.Client{},
		helmClients:         map[string]helm.Client{},
		httpClients:         map[string]*http.Client{},
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			path := filepath.Join(dir, f.Name())

			// Handle symlinks for ConfigMaps.
			ln, err := filepath.EvalSymlinks(path)
			if err == nil {
				path = ln
			}

			b, err := os.ReadFile(path)
			if err != nil {
				// Just continue if we're not able to read the 'file' as the file might be a symlink to
				// a dir when using kubernetes ConfigMaps, for example:
				//
				// drwxr-xr-x    2 root     root          4096 Oct  8 20:38 ..2020_10_08_20_38_50.434422700
				// lrwxrwxrwx    1 root     root            31 Oct  8 20:38 ..data -> ..2020_10_08_20_38_50.434422700
				continue
			}

			ac := Credentials{}

			err = json.Unmarshal(b, &ac)
			if err != nil {
				return nil, err
			}

			if ac.Name == "" {
				return nil, fmt.Errorf("no \"name\" found in artifact config file %s", path)
			}

			for _, c := range cc.artifactCredentials {
				if strings.EqualFold(ac.Name, c.Name) {
					return nil, fmt.Errorf("duplicate artifact credential listed: %s", ac.Name)
				}
			}

			// If artifact credentials is responsible for one type, generate clients as needed.
			if len(ac.Types) == 1 {
				t := ac.Types[0]
				switch t {
				case TypeGCSObject:
					opts := []option.ClientOption{option.WithScopes(storage.ScopeReadOnly)}
					if ac.JSONPath != "" {
						opts = append(opts, option.WithCredentialsFile(ac.JSONPath))
					}

					cc.gcsClients[ac.Name], err = storage.NewClient(context.Background(), opts...)
					if err != nil {
						return nil, err
					}

				case TypeGithubFile:
					var tc *http.Client

					if ac.Token != "" {
						ctx := context.Background()
						ts := oauth2.StaticTokenSource(
							&oauth2.Token{AccessToken: ac.Token},
						)
						tc = oauth2.NewClient(ctx, ts)
					}

					if ac.Enterprise {
						if ac.BaseURL == "" {
							return nil, fmt.Errorf("github file %s missing required \"baseURL\" attribute", ac.Name)
						}

						gitClient, err := github.NewEnterpriseClient(ac.BaseURL, ac.BaseURL, tc)
						if err != nil {
							return nil, err
						}

						cc.gitClients[ac.Name] = gitClient
					} else {
						gitClient := github.NewClient(tc)
						cc.gitClients[ac.Name] = gitClient
					}

				case TypeGitRepo:
					var tc *http.Client

					if ac.Token != "" {
						ctx := context.Background()
						ts := oauth2.StaticTokenSource(
							&oauth2.Token{AccessToken: ac.Token},
						)
						tc = oauth2.NewClient(ctx, ts)
					} else {
						tc = http.DefaultClient
					}

					cc.gitRepoClients[ac.Name] = tc

				case TypeHelmChart:
					if ac.Repository == "" {
						return nil, fmt.Errorf("helm chart %s missing required \"repository\" attribute", ac.Name)
					}

					helmClient := helm.NewClient(ac.Repository)

					if ac.Username != "" && ac.Password != "" {
						helmClient.WithUsernameAndPassword(ac.Username, ac.Password)
					}

					cc.helmClients[ac.Name] = helmClient

				case TypeHTTPFile:
					cc.httpClients[ac.Name] = http.DefaultClient
				}
			}

			cc.artifactCredentials = append(cc.artifactCredentials, ac)
		}
	}

	return &cc, nil
}

type credentialsController struct {
	artifactCredentials []Credentials
	httpClients         map[string]*http.Client
	helmClients         map[string]helm.Client
	gcsClients          map[string]*storage.Client
	gitClients          map[string]*github.Client
	gitRepoClients      map[string]*http.Client
}

// There might be confidential info stored in a artifacts credentials, so we need to be careful
// what we list here. In this case, only list the names and types.
func (cc *credentialsController) ListArtifactCredentialsNamesAndTypes() []Credentials {
	ac := []Credentials{}

	for _, artifaceCredentials := range cc.artifactCredentials {
		a := Credentials{
			Name:  artifaceCredentials.Name,
			Types: artifaceCredentials.Types,
		}
		ac = append(ac, a)
	}

	return ac
}

func (cc *credentialsController) HelmClientForAccountName(accountName string) (helm.Client, error) {
	if _, ok := cc.helmClients[accountName]; !ok {
		return nil, fmt.Errorf("helm account %s not found", accountName)
	}

	return cc.helmClients[accountName], nil
}

func (cc *credentialsController) HTTPClientForAccountName(accountName string) (*http.Client, error) {
	if _, ok := cc.httpClients[accountName]; !ok {
		return nil, fmt.Errorf("http account %s not found", accountName)
	}

	return cc.httpClients[accountName], nil
}

func (cc *credentialsController) GCSClientForAccountName(accountName string) (*storage.Client, error) {
	if _, ok := cc.gcsClients[accountName]; !ok {
		return nil, fmt.Errorf("gcs account %s not found", accountName)
	}

	return cc.gcsClients[accountName], nil
}

func (cc *credentialsController) GitClientForAccountName(accountName string) (*github.Client, error) {
	if _, ok := cc.gitClients[accountName]; !ok {
		return nil, fmt.Errorf("git account %s not found", accountName)
	}

	return cc.gitClients[accountName], nil
}

func (cc *credentialsController) GitRepoClientForAccountName(accountName string) (*http.Client, error) {
	if _, ok := cc.gitRepoClients[accountName]; !ok {
		return nil, fmt.Errorf("git/repo account %s not found", accountName)
	}

	return cc.gitRepoClients[accountName], nil
}
