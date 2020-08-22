package core

import (
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/billiford/go-clouddriver/pkg/sql"
	"github.com/gin-gonic/gin"
)

// I'm not sure why spinnaker needs this, but without it several necessary Spinnaker manifest stages are missing.
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
	credentials := []clouddriver.Credential{}

	providers, err := sc.ListKubernetesProviders()
	if err != nil {
		clouddriver.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	for _, provider := range providers {
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

		sca := clouddriver.Credential{
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
			// SpinnakerKindMap:        spinnakerKindMap,
			Type: "kubernetes",
		}

		if expand == "true" {
			sca.SpinnakerKindMap = spinnakerKindMap
		}
		credentials = append(credentials, sca)
	}

	c.JSON(http.StatusOK, credentials)
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
