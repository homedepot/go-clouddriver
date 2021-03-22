package kubernetes_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	"github.com/homedepot/go-clouddriver/pkg/arcade"
	"github.com/homedepot/go-clouddriver/pkg/arcade/arcadefakes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/sql"

	. "github.com/homedepot/go-clouddriver/pkg/http/core/kubernetes"
	kube "github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"github.com/homedepot/go-clouddriver/pkg/kubernetes/kubernetesfakes"
	"github.com/homedepot/go-clouddriver/pkg/sql/sqlfakes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	c                             *gin.Context
	fakeArcadeClient              *arcadefakes.FakeClient
	fakeSQLClient                 *sqlfakes.FakeClient
	fakeKubeClient                *kubernetesfakes.FakeClient
	fakeKubeController            *kubernetesfakes.FakeController
	deployManifestRequest         DeployManifestRequest
	scaleManifestRequest          ScaleManifestRequest
	cleanupArtifactsRequest       CleanupArtifactsRequest
	deleteManifestRequest         DeleteManifestRequest
	undoRolloutManifestRequest    UndoRolloutManifestRequest
	rollingRestartManifestRequest RollingRestartManifestRequest
	patchManifestRequest          PatchManifestRequest
	runJobRequest                 RunJobRequest
)

func setup() {
	gin.SetMode(gin.ReleaseMode)

	// Setup fakes.
	fakeArcadeClient = &arcadefakes.FakeClient{}

	fakeSQLClient = &sqlfakes.FakeClient{}
	fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
		Name:   "test-account",
		Host:   "http://localhost",
		CAData: "",
	}, nil)
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
	fakeKubeController.ToUnstructuredReturns(&fakeUnstructured, nil)
	fakeKubeController.SortManifestsReturns([]map[string]interface{}{
		{
			"kind":       "Pod",
			"apiVersion": "v1",
		},
	}, nil)

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req.Header.Set("X-Spinnaker-Application", "test-app")

	c, _ = gin.CreateTestContext(httptest.NewRecorder())
	c.Set(clouddriver.TaskIDKey, "test-task-id")
	c.Set(arcade.ClientInstanceKey, fakeArcadeClient)
	c.Set(kube.ControllerInstanceKey, fakeKubeController)
	c.Set(sql.ClientInstanceKey, fakeSQLClient)
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
			Cascading:          true,
			GracePeriodSeconds: &gps,
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
		Application:   "test-application",
		Manifest:      map[string]interface{}{},
		CloudProvider: "kubernetes",
		Alias:         "alias",
		Account:       "test-account",
	}
}
