package v0

import (
	"context"
	"net/http"

	clouddriver "github.com/billiford/go-clouddriver/pkg"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTask(c *gin.Context) {
	id := c.Param("id")
	manifests := []map[string]interface{}{}

	objs := cache[id]
	for _, u := range objs {
		obj := u.DeepCopyObject()
		name, err := meta.NewAccessor().Name(obj)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		gvk := obj.GetObjectKind().GroupVersionKind()

		restMapping, err := findGVR(&gvk, config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result, err := client.
			Resource(restMapping.Resource).
			Namespace(namespace).
			Get(context.TODO(), name, metav1.GetOptions{})
		// obj, err := getObject(client, *kubeconfig, o)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		manifests = append(manifests, result.Object)
	}

	ro := clouddriver.ResultObject{
		Manifests: manifests,
		CreatedArtifacts: []clouddriver.CreatedArtifact{{
			CustomKind: false,
			Location:   "",
			Metadata: struct {
				Account string "json:\"account\""
			}{
				Account: "spin-cluster-account",
			},
			Name:      "rss-site",
			Reference: "rss-site",
			Type:      "kubernetes/pod",
			Version:   "",
		},
		},
		ManifestNamesByNamespace: map[string][]string{
			"default": {"pod rss-site"},
		},
		ManifestNamesByNamespaceToRefresh: map[string][]string{
			"default": {"pod rss-site"},
		},
	}
	t := clouddriver.Task{
		ID:            id,
		ResultObjects: []clouddriver.ResultObject{ro},
		Status: clouddriver.TaskStatus{
			Complete:  true,
			Completed: true,
			Failed:    false,
			Phase:     "ORCHESTRATION",
			Retryable: false,
			Status:    "Orchestration completed.",
		},
	}
	c.JSON(http.StatusOK, t)
}
