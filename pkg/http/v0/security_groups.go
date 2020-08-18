package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SecurityGroupsResponse struct {
	Label           string           `json:"label"`
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	PropertySources []PropertySource `json:"propertySources"`
}

type PropertySource struct {
	Name   string `json:"name"`
	Source Source `json:"source"`
}

type Source struct {
	KubernetesAccounts0CacheIntervalSeconds         int    `json:"kubernetes.accounts[0].cacheIntervalSeconds"`
	KubernetesAccounts0CacheThreads                 int    `json:"kubernetes.accounts[0].cacheThreads"`
	KubernetesAccounts0CachingPolicies              string `json:"kubernetes.accounts[0].cachingPolicies"`
	KubernetesAccounts0CheckPermissionsOnStartup    bool   `json:"kubernetes.accounts[0].checkPermissionsOnStartup"`
	KubernetesAccounts0ConfigureImagePullSecrets    bool   `json:"kubernetes.accounts[0].configureImagePullSecrets"`
	KubernetesAccounts0Context                      string `json:"kubernetes.accounts[0].context"`
	KubernetesAccounts0DockerRegistries0AccountName string `json:"kubernetes.accounts[0].dockerRegistries[0].accountName"`
	KubernetesAccounts0DockerRegistries0Namespaces  string `json:"kubernetes.accounts[0].dockerRegistries[0].namespaces"`
	KubernetesAccounts0KubeconfigContents           string `json:"kubernetes.accounts[0].kubeconfigContents"`
	KubernetesAccounts0LiveManifestCalls            bool   `json:"kubernetes.accounts[0].liveManifestCalls"`
	KubernetesAccounts0Metrics                      bool   `json:"kubernetes.accounts[0].metrics"`
	KubernetesAccounts0Name                         string `json:"kubernetes.accounts[0].name"`
	KubernetesAccounts0Namespaces                   string `json:"kubernetes.accounts[0].namespaces"`
	KubernetesAccounts0OAuthScopes                  string `json:"kubernetes.accounts[0].oAuthScopes"`
	KubernetesAccounts0OmitKinds0                   string `json:"kubernetes.accounts[0].omitKinds[0]"`
	KubernetesAccounts0OmitNamespaces0              string `json:"kubernetes.accounts[0].omitNamespaces[0]"`
	KubernetesAccounts0OmitNamespaces1              string `json:"kubernetes.accounts[0].omitNamespaces[1]"`
	KubernetesAccounts0OnlySpinnakerManaged         bool   `json:"kubernetes.accounts[0].onlySpinnakerManaged"`
	KubernetesAccounts0PermissionsREAD0             string `json:"kubernetes.accounts[0].permissions.READ[0]"`
	KubernetesAccounts0PermissionsREAD1             string `json:"kubernetes.accounts[0].permissions.READ[1]"`
	KubernetesAccounts0PermissionsWRITE0            string `json:"kubernetes.accounts[0].permissions.WRITE[0]"`
	KubernetesAccounts0ProviderVersion              string `json:"kubernetes.accounts[0].providerVersion"`
	KubernetesAccounts0RequiredGroupMembership      string `json:"kubernetes.accounts[0].requiredGroupMembership"`
}

func ListSecurityGroups(c *gin.Context) {
	var empty struct{}
	// sgr := SecurityGroupsResponse{
	// 	// Label:    "securityGroups",
	// 	// Name:     "applications",
	// 	// Profiles: []string{"test2"},
	// 	// PropertySources: []PropertySource{
	// 	// 	{
	// 	// 		Name: "vualt:spinnaker",
	// 	// 		Source: Source{
	// 	// 			KubernetesAccounts0CacheIntervalSeconds:         3600,
	// 	// 			KubernetesAccounts0CacheThreads:                 4,
	// 	// 			KubernetesAccounts0CachingPolicies:              "",
	// 	// 			KubernetesAccounts0CheckPermissionsOnStartup:    false,
	// 	// 			KubernetesAccounts0ConfigureImagePullSecrets:    true,
	// 	// 			KubernetesAccounts0Context:                      "spin-cluster-account",
	// 	// 			KubernetesAccounts0DockerRegistries0AccountName: "docker-registry",
	// 	// 			KubernetesAccounts0DockerRegistries0Namespaces:  "",
	// 	// 			KubernetesAccounts0KubeconfigContents:           "",
	// 	// 			KubernetesAccounts0LiveManifestCalls:            false,
	// 	// 			KubernetesAccounts0Metrics:                      false,
	// 	// 			KubernetesAccounts0Name:                         "spin-cluster-account",
	// 	// 			KubernetesAccounts0Namespaces:                   "",
	// 	// 			KubernetesAccounts0OAuthScopes:                  "",
	// 	// 			KubernetesAccounts0OmitKinds0:                   "podPreset",
	// 	// 			KubernetesAccounts0OnlySpinnakerManaged:         true,
	// 	// 			KubernetesAccounts0PermissionsREAD0:             "gg_cloud_gcp_spinnaker_admins",
	// 	// 			KubernetesAccounts0PermissionsWRITE0:            "gg_cloud_gcp_spinnaker_admins",
	// 	// 			KubernetesAccounts0ProviderVersion:              "V2",
	// 	// 			KubernetesAccounts0RequiredGroupMembership:      "",
	// 	// 		},
	// 	// 	},
	// 	// },
	// }
	c.JSON(http.StatusOK, empty)
}

// {
//   "label": "securityGroups",
//   "name": "applications",
//   "profiles": [
//     "smoketests"
//   ],
//   "propertySources": [
//     {
//       "name": "vault:spinnaker",
//       "source": {
//         "kubernetes.accounts[0].cacheIntervalSeconds": 3600,
//         "kubernetes.accounts[0].cacheThreads": 4,
//         "kubernetes.accounts[0].cachingPolicies": "",
//         "kubernetes.accounts[0].checkPermissionsOnStartup": false,
//         "kubernetes.accounts[0].configureImagePullSecrets": true,
//         "kubernetes.accounts[0].context": "gke_github-replication-sandbox_us-central1-c_test-1-dev",
//         "kubernetes.accounts[0].dockerRegistries[0].accountName": "docker-registry",
//         "kubernetes.accounts[0].dockerRegistries[0].namespaces": "",
//         "kubernetes.accounts[0].kubeconfigContents": "apiVersion: v1\nclusters:\n- cluster:\n\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURDekNDQWZPZ0F3SUJBZ0lRSlJTSGZTc3N2SHdnZ2laTVNWbG1MakFOQmdrcWhraUc5dzBCQVFzRkFEQXYKTVMwd0t3WURWUVFERXlSbU5qaG1NMk5tWVMwMU5EazJMVFF5TkdVdFlUQTBaaTA1WmpVMVlUWTVaR001TmpJdwpIaGNOTWpBd05qRTNNakl4TnpFeVdoY05NalV3TmpFMk1qTXhOekV5V2pBdk1TMHdLd1lEVlFRREV5Um1OamhtCk0yTm1ZUzAxTkRrMkxUUXlOR1V0WVRBMFppMDVaalUxWVRZNVpHTTVOakl3Z2dFaU1BMEdDU3FHU0liM0RRRUIKQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUURsdXExN0xaY3JCVE4wa0dhdlFQeGJJN0dzbi9lUHd5eTkxZjdHNVIyaQowR1pnaGJIY0VGL0VtMUZsYUpMMjdRVWVHSFEyU01sbkkrWHNTRkt4NHROb1liSTlwQ0dZNGN0NUgrSEU1WlArCkU1YjNONlNZOEUwVEgxTFhaaDVBZ2lpelI0V3VZdkJQWlcvSFV1aFpGODc5bHdEV0lQZDU1dUt6cldQTlRPV0sKZWZxTGppcjRzRkI2QXcrSmtsNFdEajY2M3BxZEpxQ3plRzF6Qm4yVms4WU41TXBqRmtGaXppcHFZdEt4b0ZGNQoweFh3T0V5L1UyQndHaGR4QkY5djd3Ui9aUHJmdDFtWThDUVdBY3N2dU1hSXlKM25BbHdiSnFZNk9obTVMSHpwCjNJd0hzYThFMnhrYUF0K2FJSk1QbzExYkR4cDN1TXhHemN1VFhQTWU0WXlqQWdNQkFBR2pJekFoTUE0R0ExVWQKRHdFQi93UUVBd0lDQkRBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFEWgpHRStsc3pwcVo1ZEZMb2E3UDlrZ0VLaTkra1A2VXplNVpUejdiRlFZTjg3Ty9IYmxWMnlPNU9PdU0yQ1BLbGU2CktFZTZ2d2ZhcjhmMy93cUlTR25OaUtmRTdpZ0JuYzJOa1ZuY2ZMSUNXUHU4N0dQci81UEl1NHpPMmhaWnlDY00KY29Jblh6T2hSSWFQRVNmU0RQV2ErQS9vWnRiTnBlaTlKdmR3aXlGUFVndHhTS1ZFUmlHYjdXSUhpMEUwV3A4ZApaWGx3M2NQM2NETDhBQ0FlSDRXenhQS01TaVo0R3IrWHZRMzBmOWxOSHdPWngzYjF2WUJSNHZBZjllVTkwTlJ4CllKWVN4UzREeDNwejVRN0o1cWtuUVhKRUFUSyswZUI4NWt5NXcvTFVUbSticE9nMm9hK2hHV1lxK1NqcUp2NzAKUWQwRjBqUm9vYWRCa0NJaWRpQUkKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\n\n    insecure-skip-tls-verify: false\n    server: https://35.238.191.231\n  name: gke_github-replication-sandbox_us-central1-c_test-1-dev\ncontexts:\n- context:\n    cluster: gke_github-replication-sandbox_us-central1-c_test-1-dev\n    namespace: default\n    user: gke_github-replication-sandbox_us-central1-c_test-1-dev\n  name: gke_github-replication-sandbox_us-central1-c_test-1-dev\ncurrent-context: gke_github-replication-sandbox_us-central1-c_test-1-dev\nkind: Config\npreferences: {}\nusers:\n- name: gke_github-replication-sandbox_us-central1-c_test-1-dev\n  user:\n    exec:\n      apiVersion: client.authentication.k8s.io/v1beta1\n      args:\n      - /tmp/gcloud/auth_token\n      command: /bin/cat",
//         "kubernetes.accounts[0].liveManifestCalls": true,
//         "kubernetes.accounts[0].metrics": false,
//         "kubernetes.accounts[0].name": "gke_github-replication-sandbox_us-central1-c_test-1-dev",
//         "kubernetes.accounts[0].namespaces": "",
//         "kubernetes.accounts[0].oAuthScopes": "",
//         "kubernetes.accounts[0].omitKinds[0]": "podPreset",
//         "kubernetes.accounts[0].omitNamespaces[0]": "kube-public",
//         "kubernetes.accounts[0].omitNamespaces[1]": "kube-node-lease",
//         "kubernetes.accounts[0].onlySpinnakerManaged": true,
//         "kubernetes.accounts[0].permissions.READ[0]": "gg_spinnaker_admins",
//         "kubernetes.accounts[0].permissions.READ[1]": "gg_cloud_gcp_spinnaker_admins",
//         "kubernetes.accounts[0].permissions.WRITE[0]": "gg_cloud_gcp_spinnaker_admins",
//         "kubernetes.accounts[0].providerVersion": "V2",
//         "kubernetes.accounts[0].requiredGroupMembership": "",
//       }
//     }
//   ]
// }
