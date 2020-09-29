package core

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"sync"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/arcade"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

var listNamespacesTimeout = int64(5)

// I'm not sure why spinnaker needs this, but without it several necessary Spinnaker manifest stages are missing
// (I suppose this is *why* Spinnaker needs it!).
//
// Also, All accounts with this have the same kind map, so we're hardcoding it for now.
var spinnakerKindMap = map[string]string{
	"apiService":                     "unclassified",
	"clusterRole":                    "unclassified",
	"clusterRoleBinding":             "unclassified",
	"configMap":                      "configs",
	"controllerRevision":             "unclassified",
	"cronJob":                        "serverGroups",
	"customResourceDefinition":       "unclassified",
	"daemonSet":                      "serverGroups",
	"deployment":                     "serverGroupManagers",
	"event":                          "unclassified",
	"horizontalpodautoscaler":        "unclassified",
	"ingress":                        "loadBalancers",
	"job":                            "serverGroups",
	"limitRange":                     "unclassified",
	"mutatingWebhookConfiguration":   "unclassified",
	"namespace":                      "unclassified",
	"networkPolicy":                  "securityGroups",
	"persistentVolume":               "configs",
	"persistentVolumeClaim":          "configs",
	"pod":                            "instances",
	"podDisruptionBudget":            "unclassified",
	"podPreset":                      "unclassified",
	"podSecurityPolicy":              "unclassified",
	"replicaSet":                     "serverGroups",
	"role":                           "unclassified",
	"roleBinding":                    "unclassified",
	"secret":                         "configs",
	"service":                        "loadBalancers",
	"serviceAccount":                 "unclassified",
	"statefulSet":                    "serverGroups",
	"storageClass":                   "unclassified",
	"validatingWebhookConfiguration": "unclassified",
}

// List credentials for providers.
func ListCredentials(c *gin.Context) {
	expand := c.Query("expand")
	sc := sql.Instance(c)
	ac := arcade.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	credentials := []clouddriver.Credential{}

	providers, err := sc.ListKubernetesProvidersAndPermissions()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	for _, provider := range providers {
		sca := clouddriver.Credential{
			AccountType:   provider.Name,
			CloudProvider: "kubernetes",
			Environment:   provider.Name,
			Name:          provider.Name,
			Permissions: clouddriver.Permissions{
				READ:  provider.Permissions.Read,
				WRITE: provider.Permissions.Write,
			},
			PrimaryAccount:          false,
			ProviderVersion:         "v2",
			RequiredGroupMembership: []interface{}{},
			Skin:                    "v2",
			Type:                    "kubernetes",
		}

		if expand == "true" {
			sca.SpinnakerKindMap = spinnakerKindMap
		}
		credentials = append(credentials, sca)
	}

	// Only list namespaces when the 'expand' query param is set to true.
	//
	// Gate is polling the endpoint `/credentials?expand=true` once every
	// thirty seconds. Each gate instance is doing this, making the requests to get
	// all provider's namespaces a multiple of how many gate instances there are.
	if expand == "true" {
		wg := &sync.WaitGroup{}
		accountNamespacesCh := make(chan AccountNamespaces, len(providers))
		wg.Add(len(providers))

		// Get all namespaces of allowed accounts asynchronously.
		for _, provider := range providers {
			go listNamespaces(provider, wg, accountNamespacesCh, ac, kc)
		}

		wg.Wait()

		close(accountNamespacesCh)

		for an := range accountNamespacesCh {
			for i, cred := range credentials {
				if strings.EqualFold(an.Name, cred.Name) {
					cred.Namespaces = an.Namespaces
					credentials[i] = cred
				}
			}
		}
	}

	c.JSON(http.StatusOK, credentials)
}

type AccountNamespaces struct {
	Name       string
	Namespaces []string
}

func listNamespaces(provider kubernetes.Provider,
	wg *sync.WaitGroup,
	accountNamespacesCh chan AccountNamespaces,
	ac arcade.Client,
	kc kubernetes.Controller) {
	defer wg.Done()

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		log.Println("/credentials error decoding provider ca data:", err.Error())
		return
	}

	token, err := ac.Token()
	if err != nil {
		log.Println("error getting arcade token:", err.Error())
		return
	}

	config := &rest.Config{
		Host:        provider.Host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: cd,
		},
	}

	client, err := kc.NewClient(config)
	if err != nil {
		log.Println("/credentials error creating dynamic account:", err.Error())
		return
	}

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}
	// timeout listing namespaces to 5 seconds
	result, err := client.ListByGVR(gvr, metav1.ListOptions{
		TimeoutSeconds: &listNamespacesTimeout,
	})
	if err != nil {
		log.Println("/credentials error listing using kubernetes account:", err.Error())
		return
	}

	namespaces := []string{}
	for _, ns := range result.Items {
		namespaces = append(namespaces, ns.GetName())
	}
	an := AccountNamespaces{
		Name:       provider.Name,
		Namespaces: namespaces,
	}

	accountNamespacesCh <- an
}

func GetAccountCredentials(c *gin.Context) {
	sc := sql.Instance(c)
	account := c.Param("account")

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	readGroups, err := sc.ListReadGroupsByAccountName(provider.Name)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	writeGroups, err := sc.ListWriteGroupsByAccountName(provider.Name)
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	credentials := clouddriver.Credential{
		AccountType:                 provider.Name,
		ChallengeDestructiveActions: false,
		CloudProvider:               "kubernetes",
		Environment:                 provider.Name,
		Name:                        provider.Name,
		Permissions: clouddriver.Permissions{
			READ:  readGroups,
			WRITE: writeGroups,
		},
		PrimaryAccount:          false,
		ProviderVersion:         "v2",
		RequiredGroupMembership: []interface{}{},
		Skin:                    "v2",
		SpinnakerKindMap:        spinnakerKindMap,
		Type:                    "kubernetes",
	}

	c.JSON(http.StatusOK, credentials)
}
