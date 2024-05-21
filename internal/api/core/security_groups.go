package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// I'm not sure what this endpoint is supposed to do or what Spinnaker does with it,
// but here is an example of the response. I found that returning an empty object "{}"
// satisfies Deck.
//
//	{
//	  "label": "securityGroups",
//	  "name": "applications",
//	  "profiles": [
//	    "smoketests"
//	  ],
//	  "propertySources": [
//	    {
//	      "name": "vault:spinnaker",
//	      "source": {
//	        "kubernetes.accounts[0].cacheIntervalSeconds": 3600,
//	        "kubernetes.accounts[0].cacheThreads": 4,
//	        "kubernetes.accounts[0].cachingPolicies": "",
//	        "kubernetes.accounts[0].checkPermissionsOnStartup": false,
//	        "kubernetes.accounts[0].configureImagePullSecrets": true,
//	        "kubernetes.accounts[0].context": "...",
//	        "kubernetes.accounts[0].dockerRegistries[0].accountName": "docker-registry",
//	        "kubernetes.accounts[0].dockerRegistries[0].namespaces": "",
//	        "kubernetes.accounts[0].kubeconfigContents": "...",
//	        "kubernetes.accounts[0].liveManifestCalls": true,
//	        "kubernetes.accounts[0].metrics": false,
//	        "kubernetes.accounts[0].name": "...",
//	        "kubernetes.accounts[0].namespaces": "",
//	        "kubernetes.accounts[0].oAuthScopes": "",
//	        "kubernetes.accounts[0].omitKinds[0]": "podPreset",
//	        "kubernetes.accounts[0].omitNamespaces[0]": "kube-public",
//	        "kubernetes.accounts[0].omitNamespaces[1]": "kube-node-lease",
//	        "kubernetes.accounts[0].onlySpinnakerManaged": true,
//	        "kubernetes.accounts[0].permissions.READ[0]": "...",
//	        "kubernetes.accounts[0].permissions.READ[1]": "...",
//	        "kubernetes.accounts[0].permissions.WRITE[0]": "...",
//	        "kubernetes.accounts[0].providerVersion": "V2",
//	        "kubernetes.accounts[0].requiredGroupMembership": "",
//	      }
//	    }
//	  ]
//	}
func ListSecurityGroups(c *gin.Context) {
	var empty struct{}

	c.JSON(http.StatusOK, empty)
}
