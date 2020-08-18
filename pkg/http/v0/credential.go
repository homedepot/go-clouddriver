package v0

import (
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
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

// TODO lots of this is hardcoded - we need to get these from each provider.
func ListCredentials(c *gin.Context) {
	expand := c.Query("expand")
	credentials := []clouddriver.Credential{}
	sca := clouddriver.Credential{
		AccountType:                 "spin-cluster-account",
		ChallengeDestructiveActions: false,
		CloudProvider:               "kubernetes",
		Environment:                 "spin-cluster-account",
		Name:                        "spin-cluster-account",
		Permissions: clouddriver.Permissions{
			READ: []string{
				"gg_cloud_gcp_spinnaker_admins",
			},
			WRITE: []string{
				"gg_cloud_gcp_spinnaker_admins",
			},
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
	c.JSON(http.StatusOK, credentials)
}
