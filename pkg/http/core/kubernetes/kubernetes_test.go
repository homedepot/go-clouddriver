package kubernetes_test

import (
	"github.com/billiford/go-clouddriver/pkg/arcade/arcadefakes"
	. "github.com/billiford/go-clouddriver/pkg/http/core/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"github.com/billiford/go-clouddriver/pkg/kubernetes/kubernetesfakes"
	"github.com/billiford/go-clouddriver/pkg/sql/sqlfakes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	err                error
	fakeArcadeClient   *arcadefakes.FakeClient
	fakeSQLClient      *sqlfakes.FakeClient
	fakeKubeClient     *kubernetesfakes.FakeClient
	fakeKubeController *kubernetesfakes.FakeController
	actionConfig       ActionConfig
	actionHandler      ActionHandler
	action             Action
)

func setup() {
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

	ul := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							kubernetes.AnnotationSpinnakerArtifactName: "test-deployment",
							kubernetes.AnnotationSpinnakerArtifactType: "kubernetes/deployment",
							"deployment.kubernetes.io/revision":        "100",
						},
					},
				},
			},
		},
	}
	fakeKubeClient = &kubernetesfakes.FakeClient{}
	fakeKubeClient.GetReturns(&unstructured.Unstructured{Object: map[string]interface{}{}}, nil)
	fakeKubeClient.ListReturns(ul, nil)

	fakeKubeController = &kubernetesfakes.FakeController{}
	fakeKubeController.NewClientReturns(fakeKubeClient, nil)

	actionHandler = NewActionHandler()
	actionConfig = newActionConfig()
}

func newActionConfig() ActionConfig {
	return ActionConfig{
		ArcadeClient:   fakeArcadeClient,
		KubeController: fakeKubeController,
		SQLClient:      fakeSQLClient,
		ID:             "test-id",
		Application:    "test-application",
		Operation: Operation{
			DeployManifest: &DeployManifestRequest{
				EnableTraffic:     false,
				NamespaceOverride: "",
				OptionalArtifacts: nil,
				CloudProvider:     "",
				Manifests: []map[string]interface{}{
					{
						"kind":       "Pod",
						"apiVersion": "v1",
					},
				},
				TrafficManagement: struct {
					Options struct {
						EnableTraffic bool "json:\"enableTraffic\""
					} "json:\"options\""
					Enabled bool "json:\"enabled\""
				}{
					Options: struct {
						EnableTraffic bool "json:\"enableTraffic\""
					}{
						EnableTraffic: false,
					},
					Enabled: false,
				},
				Moniker: struct {
					App string "json:\"app\""
				}{
					App: "",
				},
				Source:                   "",
				Account:                  "",
				SkipExpressionEvaluation: false,
				RequiredArtifacts:        nil,
			},
			ScaleManifest: &ScaleManifestRequest{
				Replicas:      "16",
				ManifestName:  "deployment test-deployment",
				CloudProvider: "",
				Location:      "",
				User:          "",
				Account:       "",
			},
			CleanupArtifacts: &CleanupArtifactsRequest{
				Manifests: nil,
				Account:   "",
			},
			DeleteManifest: &DeleteManifestRequest{
				ManifestName:  "deployment test-deployment",
				CloudProvider: "",
				Options: struct {
					OrphanDependants   bool "json:\"orphanDependants\""
					GracePeriodSeconds int  "json:\"gracePeriodSeconds\""
				}{
					OrphanDependants:   false,
					GracePeriodSeconds: 0,
				},
				Location: "",
				User:     "",
				Account:  "",
			},
			UndoRolloutManifest: &UndoRolloutManifestRequest{
				ManifestName:  "deployment test-deployment",
				CloudProvider: "",
				Location:      "",
				User:          "",
				Account:       "",
				Revision:      "100",
			},
			RollingRestartManifest: &RollingRestartManifestRequest{
				CloudProvider: "",
				ManifestName:  "deployment test-deployment",
				Location:      "",
				User:          "",
				Account:       "",
			},
		},
	}
}
