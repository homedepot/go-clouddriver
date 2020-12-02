package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	ops "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	manifestListTimeout = int64(30)
)

func GetManifest(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	account := c.Param("account")
	namespace := c.Param("location")
	// The name of this param should really be "id" or "cluster" as it represents a Spinnaker cluster, such as "deployment my-deployment".
	// However, we have to make this path param match because of an underlying httprouter issue https://github.com/gin-gonic/gin/issues/2016.
	n := c.Param("kind")
	a := strings.Split(n, " ")
	kind := a[0]
	name := a[1]

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	token, err := ac.Token()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	result, err := client.Get(kind, name, namespace)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	app := "unknown"
	labels := result.GetLabels()

	if _, ok := labels[kubernetes.LabelKubernetesName]; ok {
		app = labels[kubernetes.LabelKubernetesName]
	}

	kmr := ops.ManifestResponse{
		Account:  account,
		Events:   []interface{}{},
		Location: namespace,
		Manifest: result.Object,
		Metrics:  []interface{}{},
		Moniker: ops.ManifestResponseMoniker{
			App:     app,
			Cluster: fmt.Sprintf("%s %s", kind, name),
		},
		Name: fmt.Sprintf("%s %s", kind, name),
		// The 'default' status of a kubernetes resource.
		Status:   kubernetes.GetStatus(kind, result.Object),
		Warnings: []interface{}{},
	}

	c.JSON(http.StatusOK, kmr)
}

func GetManifestByTarget(c *gin.Context) {
	sc := sql.Instance(c)
	kc := kubernetes.ControllerInstance(c)
	ac := arcade.Instance(c)
	account := c.Param("account")
	application := c.Param("application")
	namespace := c.Param("location")
	kind := c.Param("kind")
	cluster := c.Param("cluster")
	// Target can be newest, second_newest, oldest, largest, smallest.
	target := c.Param("target")

	// Sometimes a full kind such as MutatingWebhookConfiguration.admissionregistration.k8s.io
	// is passed in - this is the current fix for that...
	if strings.Contains(kind, ".") {
		a2 := strings.Split(kind, ".")
		kind = a2[0]
	}

	provider, err := sc.GetKubernetesProvider(account)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	cd, err := base64.StdEncoding.DecodeString(provider.CAData)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	token, err := ac.Token()
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
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
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	gvr, err := client.GVRForKind(kind)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	lo := metav1.ListOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       kind,
			APIVersion: gvr.Group + "/" + gvr.Version,
		},
		LabelSelector:  kubernetes.LabelKubernetesName + "=" + application,
		FieldSelector:  "metadata.namespace=" + namespace,
		TimeoutSeconds: &manifestListTimeout,
		// Limit:          0,
	}

	list, err := client.ListByGVR(gvr, lo)
	if err != nil {
		clouddriver.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Filter out all unassociated objects based on the moniker.spinnaker.io/cluster annotation.
	manifestFilter := kubernetes.NewManifestFilter(list.Items)
	items := manifestFilter.FilterOnCluster(cluster)
	if len(items) == 0 {
		clouddriver.Error(c, http.StatusNotFound, errors.New("no resources found for cluster "+cluster))
		return
	}

	// For now, we sort on creation timestamp to grab the manifest.
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetCreationTimestamp().String() > items[j].GetCreationTimestamp().String()
	})

	var result = items[0]

	// Target can be newest, second_newest, oldest, largest, smallest.
	// TODO fill in for largest and smallest targets.
	switch strings.ToLower(target) {
	case "newest":
		result = items[0]
	case "second_newest":
		if len(items) < 2 {
			clouddriver.Error(c, http.StatusBadRequest, errors.New("requested target \"Second Newest\" for cluster "+cluster+", but only one resource was found"))
			return
		}

		result = items[1]
	case "oldest":
		if len(items) < 2 {
			clouddriver.Error(c, http.StatusBadRequest, errors.New("requested target \"Oldest\" for cluster "+cluster+", but only one resource was found"))
			return
		}

		result = items[len(items)-1]
	default:
		clouddriver.Error(c, http.StatusNotImplemented, errors.New("requested target \""+target+"\" for cluster "+cluster+" is not supported"))
		return
	}

	kmr := ops.ManifestResponse{
		Account:  account,
		Events:   []interface{}{},
		Location: namespace,
		Manifest: result.Object,
		Metrics:  []interface{}{},
		Moniker: ops.ManifestResponseMoniker{
			App:     application,
			Cluster: cluster,
		},
		Name:     fmt.Sprintf("%s %s", kind, result.GetName()),
		Status:   kubernetes.GetStatus(kind, result.Object),
		Warnings: []interface{}{},
	}

	c.JSON(http.StatusOK, kmr)
}
