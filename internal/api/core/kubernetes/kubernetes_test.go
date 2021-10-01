package kubernetes_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/go-clouddriver/internal"
	"github.com/homedepot/go-clouddriver/internal/arcade/arcadefakes"
	"github.com/homedepot/go-clouddriver/internal/artifact"
	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"

	. "github.com/homedepot/go-clouddriver/internal/api/core/kubernetes"
	"github.com/homedepot/go-clouddriver/internal/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/internal/sql/sqlfakes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	c                             *gin.Context
	fakeArcadeClient              *arcadefakes.FakeClient
	fakeSQLClient                 *sqlfakes.FakeClient
	fakeKubeClient                *kubernetesfakes.FakeClient
	fakeKubeController            *kubernetesfakes.FakeController
	kubernetesController          *Controller
	deployManifestRequest         DeployManifestRequest
	scaleManifestRequest          ScaleManifestRequest
	cleanupArtifactsRequest       CleanupArtifactsRequest
	deleteManifestRequest         DeleteManifestRequest
	undoRolloutManifestRequest    UndoRolloutManifestRequest
	rollingRestartManifestRequest RollingRestartManifestRequest
	patchManifestRequest          PatchManifestRequest
	runJobRequest                 RunJobRequest
	clusterScopedProvider         kubernetes.Provider
	namespaceScopedProvider       kubernetes.Provider
)

func setup() {
	gin.SetMode(gin.ReleaseMode)

	var providerNamespace = "provider-namespace"

	// Setup fakes.
	fakeArcadeClient = &arcadefakes.FakeClient{}

	clusterScopedProvider = kubernetes.Provider{
		Name:      "test-account",
		Host:      "http://localhost",
		CAData:    "",
		Namespace: nil,
	}

	namespaceScopedProvider = kubernetes.Provider{
		Name:      "test-account",
		Host:      "http://localhost",
		CAData:    "",
		Namespace: &providerNamespace,
	}

	fakeSQLClient = &sqlfakes.FakeClient{}
	fakeSQLClient.GetKubernetesProviderReturns(clusterScopedProvider, nil)
	fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
		{
			AccountName: "test-account-name",
		},
	}, nil)

	fakeUnstructured := unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "test-kind",
			"apiVersion": "test-api-version",
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
					kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
					"deployment.kubernetes.io/revision":        "100",
				},
				"name": "test-name",
			},
		},
	}

	fakeUnstructuredList := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			fakeUnstructured,
		},
	}
	fakeKubeClient = &kubernetesfakes.FakeClient{}
	fakeKubeClient.GetReturns(&unstructured.Unstructured{Object: map[string]interface{}{}}, nil)
	fakeKubeClient.ListByGVRReturns(fakeUnstructuredList, nil)

	fakeKubeController = &kubernetesfakes.FakeController{}
	fakeKubeController.NewClientReturns(fakeKubeClient, nil)

	ic := &internal.Controller{
		ArcadeClient:         fakeArcadeClient,
		SQLClient:            fakeSQLClient,
		KubernetesController: fakeKubeController,
	}
	kubernetesController = &Controller{ic}

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req.Header.Set("X-Spinnaker-Application", "test-app")

	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	c.Set(clouddriver.TaskIDKey, "test-task-id")
	c.Request = req

	deployManifestRequest = newDeployManifestRequest()
	scaleManifestRequest = newScaleManifestRequest()
	cleanupArtifactsRequest = newCleanupArtifactsRequest()
	deleteManifestRequest = newDeleteManifestRequest()
	undoRolloutManifestRequest = newUndoRolloutManifestRequest()
	rollingRestartManifestRequest = newRollingRestartManifestRequest()
	patchManifestRequest = newPatchManifestRequest()
	runJobRequest = newRunJobRequest()
}

func newDeployManifestRequest() DeployManifestRequest {
	return DeployManifestRequest{
		Manifests: []map[string]interface{}{
			{
				"kind":       "Pod",
				"apiVersion": "v1",
				"metadata": map[string]interface{}{
					"namespace": "default",
					"name":      "test-name",
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "test-container-name",
							"image": "gcr.io/test-project/test-container-image",
						},
					},
				},
			},
		},
		OptionalArtifacts: []clouddriver.Artifact{
			{
				Reference: "gke-versioned-volume-config2-v004",
				Name:      "gke-versioned-volume-config2",
				Type:      artifact.TypeKubernetesConfigMap,
			},
		},
		RequiredArtifacts: []clouddriver.Artifact{
			{
				Reference: "gke-versioned-volume-config2-v004",
				Name:      "gke-versioned-volume-config2",
				Type:      artifact.TypeKubernetesConfigMap,
			},
		},
	}
}

func newScaleManifestRequest() ScaleManifestRequest {
	return ScaleManifestRequest{
		Replicas:     "16",
		ManifestName: "deployment test-deployment",
	}
}

func newCleanupArtifactsRequest() CleanupArtifactsRequest {
	return CleanupArtifactsRequest{
		Account: "test-account",
		Manifests: []map[string]interface{}{
			{
				"kind":       "Pod",
				"apiVersion": "v1",
			},
		},
	}
}

func newDeleteManifestRequest() DeleteManifestRequest {
	gps := int64(10)

	return DeleteManifestRequest{
		Location:     "test-namespace",
		Mode:         "static",
		ManifestName: "deployment test-deployment",
		Options: DeleteManifestRequestOptions{
			Cascading:          false,
			GracePeriodSeconds: &gps,
		},
		Kinds: []string{
			"deployment",
		},
		LabelSelectors: DeleteManifestRequestLabelSelectors{
			Selectors: []DeleteManifestRequestLabelSelector{
				{
					Kind:   "EQUALS",
					Key:    "key1",
					Values: []string{"key1-value1"},
				},
				{
					Kind:   "EXISTS",
					Key:    "key2",
					Values: []string{},
				},
				{
					Kind:   "NOT_CONTAINS",
					Key:    "key3",
					Values: []string{"key3-value1", "key3-value2"},
				},
			},
		},
	}
}

func newUndoRolloutManifestRequest() UndoRolloutManifestRequest {
	return UndoRolloutManifestRequest{
		ManifestName: "deployment test-deployment",
		Revision:     "100",
	}
}

func newRollingRestartManifestRequest() RollingRestartManifestRequest {
	return RollingRestartManifestRequest{
		ManifestName: "deployment test-deployment",
	}
}

func newPatchManifestRequest() PatchManifestRequest {
	return PatchManifestRequest{
		ManifestName: "deployment test-deployment",
		Options: PatchManifestRequestOptions{
			MergeStrategy: "strategic",
		},
	}
}

func newRunJobRequest() RunJobRequest {
	return RunJobRequest{
		Application: "test-application",
		Manifest: map[string]interface{}{
			"kind":       "Job",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"namespace":    "default",
				"generateName": "test-",
			},
		},
		CloudProvider: "kubernetes",
		Alias:         "alias",
		Account:       "test-account",
	}
}
