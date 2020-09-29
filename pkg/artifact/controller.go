package artifact

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/billiford/go-clouddriver/pkg/helm"
	// "github.com/google/go-github/v32/github"
)

type Type string

const (
	TypeHelmChart               Type = "helm/chart"
	TypeGitRepo                 Type = "git/repo"
	TypeFront50PipelineTemplate Type = "front50/pipelineTemplate"
	TypeEmbeddedBase64          Type = "embedded/base64"
	TypeCustomerObject          Type = "custom/object"
	TypeGCSObject               Type = "gcs/object"
	TypeDockerImage             Type = "docker/image"
	TypeKubernetesConfigMap     Type = "kubernetes/configMap"
	TypeKubernetesDeployment    Type = "kubernetes/deployment"
	TypeKubernetesReplicaSet    Type = "kubernetes/replicaSet"
	TypeKubernetesSecret        Type = "kubernetes/secret"
	TypeGithubFile              Type = "github/file"
)

type CredentialsController interface {
	ListArtifactCredentialsNamesAndTypes() []Credentials
	HelmClientForAccountName(string) (helm.Client, error)
	// GitClientForAccountName(string) (github.Client, error)
}

type Credentials struct {
	Name  string `json:"name"`
	Types []Type `json:"types"`
	URL   string `json:"url,omitempty"`
}

var (
	defaultConfigDir = "/opt/spinnaker/artifacts/config"
)

func NewDefaultCredentialsController() (CredentialsController, error) {
	return NewCredentialsController(defaultConfigDir)
}

func NewCredentialsController(dir string) (CredentialsController, error) {
	cc := credentialsController{
		artifactCredentials: []Credentials{},
		helmClients:         map[string]helm.Client{},
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			b, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}

			ac := Credentials{}

			err = json.Unmarshal(b, &ac)
			if err != nil {
				return nil, err
			}

			if ac.Name == "" {
				return nil, fmt.Errorf("no \"name\" found in artifact config file %s", filepath.Join(dir, f.Name()))
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
				case TypeHelmChart:
					if ac.URL == "" {
						return nil, fmt.Errorf("helm chart %s missing required \"url\" attribute", ac.Name)
					}

					helmClient := helm.NewClient(ac.URL)
					cc.helmClients[ac.Name] = helmClient
				}
			}

			cc.artifactCredentials = append(cc.artifactCredentials, ac)
		}
	}

	return &cc, nil
}

type credentialsController struct {
	artifactCredentials []Credentials
	helmClients         map[string]helm.Client
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
